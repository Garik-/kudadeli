package web

import (
	"encoding/json"
	"io"
	"kudadeli/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func expensesHandler(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h := w.Header()

		lastModified, err := db.LatestUpdatedAt(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "db.LatestUpdatedAt:", "error", err)
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

		lastModified = expenses.LatestUpdatedAt()

		h.Set("Content-Type", "application/json; charset=utf-8")
		h.Set("Cache-Control", "private, must-revalidate")
		h.Set("Last-Modified", lastModified.UTC().Format(http.TimeFormat))

		w.WriteHeader(http.StatusOK)

		if len(expenses) == 0 {
			_, err = io.WriteString(w, "[]")
			if err != nil {
				slog.ErrorContext(ctx, "writeString", "error", err)
				writeError(w, err.Error())
			}

			return
		}

		if err := json.NewEncoder(w).Encode(expenses); err != nil {
			slog.ErrorContext(ctx, "writeString", "error", err)
			writeError(w, err.Error())
		}
	}
}

type updateExpenseCategoryRequest struct {
	Category byte `json:"category" validate:"required,gt=0"`
}

func updateExpenseCategoryHandler(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()

		ctx := r.Context()
		slog.DebugContext(ctx, "update category handler")

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			writeErrorWithCode(w, err.Error(), http.StatusBadRequest)

			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var req updateExpenseCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorWithCode(w, err.Error(), http.StatusBadRequest)

			return
		}

		category := model.Category(req.Category)

		if !category.IsValid() {
			writeErrorWithCode(w, "category ID must be positive byte", http.StatusBadRequest)

			return
		}

		if err := db.UpdateCategory(ctx, id, model.Category(req.Category)); err != nil {
			slog.ErrorContext(ctx, "UpdateCategory", "error", err)
			writeError(w, "failed to update category")

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
