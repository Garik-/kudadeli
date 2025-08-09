package web

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

func authMiddleware(token string, allowedUsers []int64, expIn time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			ctx := r.Context()

			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 {
				slog.ErrorContext(ctx, "malformed Authorization header: expected format '<type> <token>'")
				writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			authType := authParts[0]
			authData := authParts[1]

			if authType != "tma" {
				slog.ErrorContext(ctx, "invalid authorization type: expected 'tma'")
				writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			// Validate init data (valid 1 hour)
			if err := initdata.Validate(authData, token, expIn); err != nil {
				slog.ErrorContext(ctx, "initdata.Validate", "error", err.Error())
				writeErrorWithCode(w, err.Error(), http.StatusUnauthorized)

				return
			}

			if len(allowedUsers) > 0 {
				// Parse init data
				initData, err := initdata.Parse(authData)
				if err != nil {
					slog.ErrorContext(ctx, "initdata.Parse", "error", err.Error())
					writeError(w, err.Error())

					return
				}

				if !slices.Contains(allowedUsers, initData.User.ID) {
					slog.WarnContext(ctx, "unauthorized: user ID is not in the allowed users", "user_id", initData.User.ID)
					writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
