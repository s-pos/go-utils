package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

const (
	BEARER_TOKEN = "Bearer"
	CTX_KEY_USER = "user_id"
)

func (c *client) Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var (
			ctx  = req.Context()
			sess sessionAuth
			auth = req.Header.Get("Authorization")
			err  error
		)

		// check header authorization is null or not
		if reflect.ValueOf(auth).IsZero() {
			err = errors.New("header authorization not found")
			c.errorResponse(ctx, w, http.StatusForbidden, string(HeaderAuthorizationNotFound), resMessage[HeaderAuthorizationNotFound], resReason[HeaderAuthorizationNotFound], err)
			return
		}

		token := strings.Split(auth, " ")
		if len(token) < 2 {
			err = fmt.Errorf("len of token %s", token)
			c.errorResponse(ctx, w, http.StatusUnauthorized, string(SessionError), resMessage[SessionError], resReason[SessionError], err)
			return
		}
		if token[0] != BEARER_TOKEN {
			err = fmt.Errorf("token type not same. from header %s", token[0])
			c.errorResponse(ctx, w, http.StatusUnauthorized, string(SessionTokenTypeWrong), resMessage[SessionTokenTypeWrong], resReason[SessionTokenTypeWrong], err)
			return
		}

		result, err := c.redis.Get(ctx, token[1]).Result()
		if err != nil {
			err = fmt.Errorf("error get token from redis %v", err)
			c.errorResponse(ctx, w, http.StatusUnauthorized, string(SessionNotFound), resMessage[SessionNotFound], resReason[SessionNotFound], err)
			return
		}

		err = json.Unmarshal([]byte(result), &sess)
		if err != nil {
			err = fmt.Errorf("error when unmarshal token to interface. %v", err)
			c.errorResponse(ctx, w, http.StatusInternalServerError, string(SessionError), resMessage[SessionError], resReason[SessionError], err)
			return
		}

		// store to context
		ctx = context.WithValue(
			ctx,
			CTX_KEY_USER,
			sess.ID,
		)
		// store to header
		req.Header.Set("x-sess-token", result)
		// continue to controller handler
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
