package web

import (
	"context"
	"encoding/json"
	"kudadeli/model"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

const (
	readHeaderTimeout = 2 * time.Second
)

type Database interface {
	List(ctx context.Context, limit int) (model.Expenses, error)
	LatestUpdatedAt(ctx context.Context) (time.Time, error)
	UpdateCategory(ctx context.Context, expenseID model.ExpenseID, category model.Category) error
}

func newServer(ctx context.Context, addr string) *http.Server {
	return &http.Server{
		ReadHeaderTimeout: readHeaderTimeout,
		Addr:              addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}

func writeError(w http.ResponseWriter, err string) {
	writeErrorWithCode(w, err, http.StatusInternalServerError)
}

func writeErrorWithCode(w http.ResponseWriter, err string, statusCode int) {
	h := w.Header()

	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": err,
	}); err != nil {
		slog.Error("json encode: %w", "error", err)
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

func New(ctx context.Context, addr string, allowedOrigins []string, db Database) (*http.Server, error) {
	fs := http.FileServer(http.FS(publicFiles))

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r := chi.NewRouter()
	r.Handle("/", publicHandler(fs))
	r.NotFound(publicHandler(fs))

	r.Route("/v1", func(v1 chi.Router) {
		v1.Use(c.Handler)
		v1.Use(middleware.Timeout(2 * time.Second))

		v1.Get("/expenses", expensesHandler(db))
		v1.Put("/expenses/{id}/category", updateExpenseCategoryHandler(db))
		v1.Get("/categories", categoriesHandler())
	})

	srv := newServer(ctx, addr)
	srv.Handler = r

	return srv, nil
}
