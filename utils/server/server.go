package server

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Wrapper wraping the routing server with cors handler
func Wrapper(route *echo.Echo) error {
	// cors condition
	origins := strings.Split(os.Getenv("CORS_ORIGIN"), ",")
	headers := strings.Split(os.Getenv("CORS_HEADERS"), ",")
	log.Println(origins, headers)
	route.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	route.HideBanner = true

	// check liveness and readiness connection
	hc := NewHealthCheck()
	health := hc.HealthCheck()
	// route group for liveness, readines and prometheus
	rg := route.Group("/live")
	// route for healthcheck
	rg.GET("/status", echo.WrapHandler(health.Handler()))

	// route for metrics prometheus
	rg.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	return route.Start(port)
}
