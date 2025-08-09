package web

import (
	"net/http"
	"slices"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// TODO: add error logs
func authMiddleware(token string, allowedUsers []int64, expIn time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			authParts := strings.Split(authHeader, " ")
			if len(authParts) != 2 {
				writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			authType := authParts[0]
			authData := authParts[1]

			if authType != "tma" {
				writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			// Validate init data (valid 1 hour)
			if err := initdata.Validate(authData, token, expIn); err != nil {
				writeErrorWithCode(w, err.Error(), http.StatusUnauthorized)

				return
			}

			if len(allowedUsers) > 0 {
				// Parse init data
				initData, err := initdata.Parse(authData)
				if err != nil {
					writeError(w, err.Error())

					return
				}

				if !slices.Contains(allowedUsers, initData.User.ID) {
					writeErrorWithCode(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
