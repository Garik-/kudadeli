package web

import (
	"encoding/json"
	"fmt"
	"kudadeli/model"
	"log/slog"
	"net/http"
	"time"
)

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
