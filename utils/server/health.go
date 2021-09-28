package server

import (
	"os"
	"strconv"
	"time"

	"github.com/hellofresh/health-go/v4"
	pg "github.com/hellofresh/health-go/v4/checks/postgres"
)

type healthCheck struct {
	health *health.Health
}

type HealthCheck interface {
	HealthCheck() *health.Health
}

func NewHealthCheck() HealthCheck {
	health, _ := health.New()
	return &healthCheck{health: health}
}

func (h *healthCheck) HealthCheck() *health.Health {
	db, _ := strconv.ParseBool(os.Getenv("DB_ENABLED"))

	// check postgres connection
	if db {
		h.health.Register(health.Config{
			Name:      "postgres",
			Timeout:   3 * time.Second,
			SkipOnErr: false,
			Check: pg.New(pg.Config{
				DSN: os.Getenv("DB_URL"),
			}),
		})
	}

	return h.health
}
