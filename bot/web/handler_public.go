package web

import (
	"embed"
	"errors"
	"fmt"
	iofs "io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

//go:embed public
var publicFiles embed.FS

func fileExists(fs embed.FS, path string) bool {
	_, err := fs.Open(path)

	return !errors.Is(err, iofs.ErrNotExist)
}

func publicHandler(fileServer http.Handler) http.HandlerFunc {
	etagCache := make(map[string]string)
	version := time.Now().Unix()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		slog.DebugContext(r.Context(), r.URL.Path)

		r.URL.Path = "public" + r.URL.Path

		if r.URL.Path == "public/" || fileExists(publicFiles, r.URL.Path) {
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

			switch {
			case strings.Contains(r.URL.Path, "/assets/"):
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

			case r.URL.Path == "public/" || strings.HasSuffix(r.URL.Path, "index.html"):
				w.Header().Set("Cache-Control", "no-cache")

			default:
				w.Header().Set("Cache-Control", "public, max-age=3600")
			}
		}

		fileServer.ServeHTTP(w, r)
	}
}
