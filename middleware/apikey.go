package middleware

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
)

func (c *client) APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			ctx    = req.Context()
			key    = req.Header.Get("x-api-key") // from header request
			apiKey = os.Getenv("API_KEY")        // from environment
			err    error
		)

		if reflect.ValueOf(key).IsZero() {
			err = fmt.Errorf("api key not found in header")

			c.errorResponse(ctx, w, http.StatusForbidden, string(APIKeyNotFound), resMessage[APIKeyNotFound], resReason[APIKeyNotFound], err)
			return
		}

		if key != apiKey {
			err = fmt.Errorf("api key not same with environment. header %s. from environment %s", key, apiKey)

			c.errorResponse(ctx, w, http.StatusForbidden, string(APIKeyInvalid), resMessage[APIKeyInvalid], resReason[APIKeyInvalid], err)
			return
		}

		next.ServeHTTP(w, req)
	})
}
