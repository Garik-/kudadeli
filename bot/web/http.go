package web

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	iofs "io/fs"
	"kudadeli/model"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
)

const (
	readHeaderTimeout = 2 * time.Second
)

type Database interface {
	List(ctx context.Context) ([]model.Expense, error)
}

//go:embed public
var publicFiles embed.FS

func newServer(ctx context.Context, addr string) *http.Server {
	return &http.Server{
		ReadHeaderTimeout: readHeaderTimeout,
		Addr:              addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}

func fileExists(fs embed.FS, path string) bool {
	_, err := fs.Open(path)

	return !errors.Is(err, iofs.ErrNotExist)
}

func mainHandler(fileServer http.Handler) http.HandlerFunc {
	etagCache := make(map[string]string)
	version := time.Now().Unix()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		r.URL.Path = "public" + r.URL.Path

		if fileExists(publicFiles, r.URL.Path) {
			if etag, ok := etagCache[r.URL.Path]; ok {
				w.Header().Set("ETag", etag)

				if match := r.Header.Get("If-None-Match"); match == etag {
					slog.DebugContext(r.Context(), "static file not modified")
					w.WriteHeader(http.StatusNotModified)

					return
				}
			} else {
				etag := fmt.Sprintf(`W/"%s-%d"`, r.URL.Path, version)
				w.Header().Set("ETag", etag)
				etagCache[r.URL.Path] = etag
			}

			w.Header().Set("Cache-Control", "public, max-age=3600")
		}

		fileServer.ServeHTTP(w, r)
	}
}

func writeError(w http.ResponseWriter, err string) {
	h := w.Header()

	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": err,
	}); err != nil {
		slog.Error("json encode: %w", "error", err)
	}
}

func latestUpdatedAt(expenses []model.Expense) time.Time {
	var latest time.Time
	for i := range expenses {
		if expenses[i].UpdatedAt.After(latest) {
			latest = expenses[i].UpdatedAt
		}
	}

	return latest
}

func expensesHandler(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		ctx := r.Context()

		expenses, err := db.List(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "db.List: %w", "error", err)
			writeError(w, err.Error())

			return
		}

		h := w.Header()
		lastModified := latestUpdatedAt(expenses)

		if !lastModified.IsZero() {
			ifModifiedSince := r.Header.Get("If-Modified-Since")

			if t, err := http.ParseTime(ifModifiedSince); err == nil && !lastModified.After(t) {
				slog.DebugContext(ctx, "expenses not modified")
				w.WriteHeader(http.StatusNotModified)

				return
			}

			h.Set("Last-Modified", lastModified.Format(http.TimeFormat))
		}

		h.Set("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		if len(expenses) == 0 {
			_, err = io.WriteString(w, "[]")
			if err != nil {
				slog.ErrorContext(ctx, "writeString: %w", "error", err)
				writeError(w, err.Error())
			}

			return
		}

		if err := json.NewEncoder(w).Encode(expenses); err != nil {
			slog.ErrorContext(ctx, "writeString: %w", "error", err)
			writeError(w, err.Error())
		}
	}
}

func New(ctx context.Context, addr string, allowedOrigins []string, db Database) (*http.Server, error) {
	fs := http.FileServer(http.FS(publicFiles))

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	mux := http.NewServeMux()

	srv := newServer(ctx, addr)
	srv.Handler = mux

	mux.Handle("/", mainHandler(fs))
	mux.Handle("/v1/expenses", c.Handler(expensesHandler(db)))

	return srv, nil
}
