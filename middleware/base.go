package middleware

import (
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type client struct {
	redis *redis.Client
	log   *logrus.Logger
	loc   *time.Location
}

// Clients create contract middleware packages
type Clients interface {
	// APIKey Check API Key header for access restricted endpoint
	APIKey(next http.Handler) http.Handler

	// Logger Record custom log using Logrus library
	Logger(next http.Handler) http.Handler

	// Session Check session authrization for access the restrict endpoint
	Session(next http.Handler) http.Handler

	// CacheSpecific Get data from redis with specific key. need endpoint and user_id from user
	CacheSpecific(next http.Handler) http.Handler

	// CacheGlobal Get data from redis with key of redis only endpoint of services
	CacheGlobal(next http.Handler) http.Handler
}

// NewMiddleware will create an object that represent clients interface
func NewMiddleware(redis *redis.Client, log *logrus.Logger, loc *time.Location) Clients {
	return &client{
		redis: redis,
		log:   log,
		loc:   loc,
	}
}
