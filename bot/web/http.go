package web

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"kudadeli/model"
	"log/slog"
	"net"
	"net/http"
	"time"
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

func mainHandler(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		r.URL.Path = "public" + r.URL.Path
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

func New(ctx context.Context, addr string, db Database) (*http.Server, error) {
	fs := http.FileServer(http.FS(publicFiles))

	mux := http.NewServeMux()

	srv := newServer(ctx, addr)
	srv.Handler = mux

	mux.Handle("/", mainHandler(fs))
	mux.Handle("/v1/expenses", expensesHandler(db))

	return srv, nil
}
