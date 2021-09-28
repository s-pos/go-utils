package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/s-pos/go-utils/adapter"
	"github.com/s-pos/go-utils/utils/monitoring"
	"github.com/sirupsen/logrus"
)

// all variables config
const (
	// any error message
	ErrTimezone = "Sorry. server configuration not available time Asia/Jakarta"
	LayoutDate  = "2006-01-02 15:04:05"
	tzLocation  = "Asia/Jakarta"
)

var ServiceName string

// Load any configuration like open connection database, open connection redis, monitoring, e.t.c
func Load(serviceName string) {
	// make a servicename to lower first then if any spacebar, replace to dash (-)
	serviceName = strings.ReplaceAll(strings.ToLower(serviceName), " ", "-")
	// load prometheus
	monitoring.NewPrometheus(serviceName)
	ServiceName = serviceName

	// check database is enabled or not
	dbStatus, _ := strconv.ParseBool(os.Getenv("DB_ENABLED"))
	if dbStatus {
		adapter.LoadDatabase()
	}

	// check redis is enabled or not
	redisStatus, _ := strconv.ParseBool(os.Getenv("REDIS_ENABLED"))
	if redisStatus {
		adapter.LoadRedis()
	}
}

// Timezone load timezone area
func Timezone() *time.Location {
	loc, err := time.LoadLocation(tzLocation)
	if err != nil {
		panic(ErrTimezone)
	}
	return loc
}

// Logrus loggrus formatter json
func Logrus() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: LayoutDate,
	})
	return log
}

func GetServiceName() string {
	return ServiceName
}
