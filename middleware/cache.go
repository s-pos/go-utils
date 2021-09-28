package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/s-pos/go-utils/logger"
)

func (c *client) CacheSpecific(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			ctx         = req.Context()
			endpoint    = req.URL.Path
			token       = req.Header.Get("x-sess-token")
			key         string
			userSession sessionAuth
		)

		err := json.Unmarshal([]byte(token), &userSession)
		if err != nil {
			logger.Messagef("failed unmasrhal %v", err).To(ctx)
			// skip, let them in to endpoint handler because failed marshal
			next.ServeHTTP(w, req)
			return
		}

		// combine endpoint with user_id from session header
		key = fmt.Sprintf("%s|%d", endpoint, userSession.ID)
		// convert to base64 encode
		key = base64.StdEncoding.EncodeToString([]byte(key))

		result, err := c.redis.Get(ctx, key).Result()
		if err != nil {
			// just skip and let them in to endpoint handler for make new request data
			next.ServeHTTP(w, req)
			return
		}

		var data interface{}
		err = json.Unmarshal([]byte(result), &data)
		if err != nil {
			// just skip and let them in to endpoint handler for make new request data
			next.ServeHTTP(w, req)
			return
		}

		c.successResponse(ctx, w, http.StatusOK, string(CacheSuccess), resMessage[CacheSuccess], resReason[CacheSuccess], data)
	})
}

func (c *client) CacheGlobal(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			ctx      = req.Context()
			endpoint = req.URL.Path
			key      = endpoint
		)

		// convert to base64 encode
		key = base64.StdEncoding.EncodeToString([]byte(key))

		result, err := c.redis.Get(ctx, key).Result()
		if err != nil {
			// just skip and let them in to endpoint handler for make new request data
			next.ServeHTTP(w, req)
			return
		}

		var data interface{}
		err = json.Unmarshal([]byte(result), &data)
		if err != nil {
			// just skip and let them in to endpoint handler for make new request data
			next.ServeHTTP(w, req)
			return
		}

		c.successResponse(ctx, w, http.StatusOK, string(CacheSuccess), resMessage[CacheSuccess], resReason[CacheSuccess], data)
	})
}
