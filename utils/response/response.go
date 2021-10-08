package response

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/s-pos/go-utils/logger"
)

var CtxUserId = "user_id"

type response struct {
	Status struct {
		Error   bool   `json:"error"`
		Code    string `json:"code"`
		Message string `json:"message"`
		Reason  string `json:"reason,omitempty"`
	} `json:"status"`

	StatusCode int         `json:"-"`
	Data       interface{} `json:"data,omitempty"`

	MandatoryFields []logger.MandatoryField `json:"required_fields,omitempty"`

	// for caching
	Cache         bool          `json:"-"`
	CacheSpesific bool          `json:"-"`
	duration      time.Duration `json:"-"`
	redisClient   *redis.Client `json:"-"`
}

type Output interface {
	Write(c echo.Context) error

	Middleware(w http.ResponseWriter)
}

func (r *response) Middleware(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	json.NewEncoder(w).Encode(r)
}

func (r *response) Write(c echo.Context) error {
	if r.Cache {
		var (
			req      = c.Request()
			ctx      = req.Context()
			endpoint = req.URL.Path
			key      = endpoint
		)

		userId, ok := ctx.Value(CtxUserId).(int)
		if ok && r.CacheSpesific {
			key = fmt.Sprintf("%s|%d", key, userId)
		}

		data, err := json.Marshal(r.Data)
		if err == nil {
			r.redisClient.Set(ctx, key, string(data), r.duration)
		}
	}

	return c.JSON(r.StatusCode, r)
}

// return error response with json
func Errors(ctx context.Context, responseStatusCode int, statusCode, statusMessage, statusReason string, err error) Output {
	var errLocation string
	if fname, _, line, ok := runtime.Caller(1); ok {
		errLocation = fmt.Sprintf("[%s:%d]", runtime.FuncForPC(fname).Name(), line)
	}

	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = true
	res.Status.Code = statusCode
	res.Status.Message = statusMessage
	res.Status.Reason = statusReason
	res.MandatoryFields = logger.GetMandatoryFields(ctx)

	logger.Response(ctx, responseStatusCode, res, errLocation, err)
	logger.ResponseMessage(ctx, statusMessage)

	return &res
}

// return error response with json
func ErrorsWithData(ctx context.Context, responseStatusCode int, statusCode, statusMessage, statusReason string, data interface{}, err error) Output {
	var errLocation string
	if fname, _, line, ok := runtime.Caller(1); ok {
		errLocation = fmt.Sprintf("[%s:%d]", runtime.FuncForPC(fname).Name(), line)
	}

	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = true
	res.Status.Code = statusCode
	res.Status.Message = statusMessage
	res.Status.Reason = statusReason
	res.MandatoryFields = logger.GetMandatoryFields(ctx)

	logger.Response(ctx, responseStatusCode, res, errLocation, err)
	logger.ResponseMessage(ctx, statusMessage)

	res.Data = data

	return &res
}

// return success response with json
func Success(ctx context.Context, responseStatusCode int, statusCode, statusMessage string, data interface{}) Output {
	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = false
	res.Status.Code = statusCode
	res.Status.Message = statusMessage

	// send to logger before set data to response
	logger.Response(ctx, responseStatusCode, res, nil, nil)
	logger.ResponseMessage(ctx, statusMessage)

	res.Data = data

	return &res
}

// return success response with cache endpoint
func SuccessWithCache(ctx context.Context, responseStatusCode int, statusCode, statusMessage string, data interface{}, redis *redis.Client, duration time.Duration) Output {
	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = false
	res.Status.Code = statusCode
	res.Status.Message = statusMessage

	// send to logger before set data to response
	logger.Response(ctx, responseStatusCode, res, nil, nil)
	logger.ResponseMessage(ctx, statusMessage)

	res.Data = data
	res.Cache = true
	res.duration = duration
	res.redisClient = redis

	return &res
}

// return success response with cache endpoint + userId
func SuccessWithCacheSpesific(ctx context.Context, responseStatusCode int, statusCode, statusMessage string, data interface{}, redis *redis.Client, duration time.Duration) Output {
	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = false
	res.Status.Code = statusCode
	res.Status.Message = statusMessage

	// send to logger before set data to response
	logger.Response(ctx, responseStatusCode, res, nil, nil)
	logger.ResponseMessage(ctx, statusMessage)

	res.Data = data
	res.Cache = true
	res.CacheSpesific = true
	res.duration = duration
	res.redisClient = redis

	return &res
}

// return success response with json
func SuccessWithReason(ctx context.Context, responseStatusCode int, statusCode, statusMessage, statusReason string, data interface{}) Output {
	res := response{}
	res.StatusCode = responseStatusCode
	res.Status.Error = false
	res.Status.Code = statusCode
	res.Status.Message = statusMessage
	res.Status.Reason = statusReason

	// send to logger before set data to response
	logger.Response(ctx, responseStatusCode, res, nil, nil)
	logger.ResponseMessage(ctx, statusMessage)

	res.Data = data

	return &res
}
