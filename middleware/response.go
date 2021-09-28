package middleware

import (
	"context"
	"net/http"

	"github.com/s-pos/go-utils/utils/response"
)

type Code string

const (
	NotUser                     Code = "909980"
	NotAdmin                    Code = "909981"
	NotProductOwner             Code = "909982"
	CacheSuccess                Code = "909990"
	HeaderAuthorizationNotFound Code = "909991"
	SessionNotFound             Code = "909992" // session not found or session expired
	SessionTokenTypeWrong       Code = "909993"
	SessionLength               Code = "909994"
	APIKeyNotFound              Code = "909997"
	APIKeyInvalid               Code = "909998"

	// 	for global
	ErrorMarshal Code = "909996"
	SessionError Code = "909999"
)

var (
	resMessage = map[Code]string{
		NotUser:                     "session.roles.user",
		NotAdmin:                    "session.roles.admin",
		NotProductOwner:             "session.roles.product_owner",
		CacheSuccess:                "cache.success",
		HeaderAuthorizationNotFound: "header.authorization.required",
		SessionTokenTypeWrong:       "token.invalid",
		SessionLength:               "header.authorization.length",
		APIKeyNotFound:              "header.apikey.required",
		APIKeyInvalid:               "header.apikey.invalid",
		SessionError:                "session.invalid",
		ErrorMarshal:                "system.marshal.error",
	}

	resReason = map[Code]string{
		NotUser:                     "Hanya pengguna yang bisa akses halaman ini",
		NotAdmin:                    "Hanya admin yang bisa akses halaman ini",
		NotProductOwner:             "Hanya Product Owner yang bisa akses halaman ini",
		CacheSuccess:                "Sukses mengambil data",
		HeaderAuthorizationNotFound: "Dibutuhkan Authorization untuk mengakses halaman iini",
		SessionNotFound:             "Session sudah tidak berlaku atau tidak ditemukan",
		SessionTokenTypeWrong:       "Invalid token",
		APIKeyNotFound:              "Membutuhkan API Key untuk mengakses halaman ini",
		APIKeyInvalid:               "API Key invalid",
		SessionError:                "Session invalid",
		ErrorMarshal:                "Terjadi kesalahan pada server, silahkan coba beberapa saat lagi",
	}
)

func (c *client) errorResponse(ctx context.Context, w http.ResponseWriter, statusCode int, systemCode, message, reason string, err error) {
	response.Errors(ctx, statusCode, systemCode, message, reason, err).Middleware(w)
}

func (c *client) successResponse(ctx context.Context, w http.ResponseWriter, statusCode int, systemCode, message, reason string, data interface{}) {
	response.SuccessWithReason(ctx, statusCode, systemCode, message, reason, data).Middleware(w)
}
