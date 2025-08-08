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

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

const (
	readHeaderTimeout = 2 * time.Second
)

type Database interface {
	List(ctx context.Context, limit int) ([]model.Expense, error)
	LatestUpdatedAt(ctx context.Context) (time.Time, error)
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

type categoryJSON struct {
	ID   byte   `json:"id"`
	Name string `json:"name"`
}

func categoriesHandler() http.HandlerFunc {
	etag := fmt.Sprintf(`"W/%d"`, time.Now().Unix())

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		h := w.Header()
		h.Set("ETag", etag)

		if match := r.Header.Get("If-None-Match"); match == etag {
			slog.DebugContext(ctx, "category not modified")
			w.WriteHeader(http.StatusNotModified)

			return
		}

		categories := model.Categories()

		var jsonData = make([]categoryJSON, len(categories))

		for i := range categories {
			jsonData[i] = categoryJSON{
				ID:   byte(categories[i]),
				Name: categories[i].String(),
			}
		}

		h.Set("Content-Type", "application/json; charset=utf-8")
		h.Set("Cache-Control", "public, max-age=3600, must-revalidate")

		if err := json.NewEncoder(w).Encode(jsonData); err != nil {
			slog.ErrorContext(ctx, "writeString: %w", "error", err)
			writeError(w, err.Error())
		}
	}
}

func updateExpenseCategoryHandler(_ Database) http.HandlerFunc {
	return func(_ http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.DebugContext(ctx, "update category handler")
	}
}

func isNotModified(clientSince, lastModified time.Time) bool {
	if lastModified.IsZero() {
		return false
	}

	clientSinceUTC := clientSince.UTC()
	lastModifiedUTC := lastModified.UTC()

	return !lastModifiedUTC.After(clientSinceUTC)
}

func expensesHandler(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h := w.Header()

		lastModified, err := db.LatestUpdatedAt(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "db.LatestUpdatedAt: %w", "error", err)
			writeError(w, err.Error())

			return
		}

		// Проверяем If-Modified-Since
		ifModifiedSince := r.Header.Get("If-Modified-Since")
		if ifModifiedSince != "" {
			clientSince, err := http.ParseTime(ifModifiedSince)
			if err != nil {
				slog.ErrorContext(ctx, "failed to parse If-Modified-Since", "error", err)
			} else if isNotModified(clientSince, lastModified) {
				slog.DebugContext(ctx, "expenses not modified", "lastModified", lastModified, "clientSince", clientSince)
				w.WriteHeader(http.StatusNotModified)

				return
			}
		}

		slog.DebugContext(ctx, "expenses query")

		expenses, err := db.List(ctx, -1) // -1 means no limit
		if err != nil {
			slog.ErrorContext(ctx, "db.List: %w", "error", err)
			writeError(w, err.Error())

			return
		}

		lastModified = latestUpdatedAt(expenses)

		h.Set("Content-Type", "application/json; charset=utf-8")
		h.Set("Cache-Control", "private, must-revalidate")
		h.Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))

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

	r := chi.NewRouter()
	r.Handle("/", mainHandler(fs))
	r.Route("/v1", func(v1 chi.Router) {
		v1.Use(c.Handler)

		v1.Get("/expenses", expensesHandler(db))
		v1.Get("/categories", categoriesHandler())
		v1.Put("/expenses/{id}/category", updateExpenseCategoryHandler(db))
	})

	srv := newServer(ctx, addr)
	srv.Handler = r

	return srv, nil
}
