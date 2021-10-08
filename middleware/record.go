package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/s-pos/go-utils/logger"
	"github.com/s-pos/go-utils/utils/monitoring"
	"github.com/sirupsen/logrus"
)

var (
	colorReset   = "\033[0m"
	colorDanger  = "\033[31m"
	colorSuccess = "\033[32m"
	colorWarning = "\033[33m"
)

// Logger using for record any request, response and some logmessages to stdout terminal
func (c *client) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now().In(c.loc)
		req, data := logger.Initialize(r, start)

		next.ServeHTTP(rw, req)

		data.Finalize(req.Context())
		c.write(&data)
	})
}

func (c *client) write(dl *logger.DataLogger) {
	var (
		level            logrus.Level
		elasticStatus, _ = strconv.ParseBool(os.Getenv("ELASTIC_ENABLED"))
		errChan          = make(chan error, 1)
	)
	if elasticStatus {
		go func() {
			elastic := logger.NewElastic()
			elastic.SendDataToElastic(dl, errChan)
		}()
	}

	monitoring.Prometheus().Record(dl.StatusCode, dl.RequestMethod, dl.Endpoint, dl.ResponseMessage, time.Since(dl.TimeStart))
	if elasticStatus {
		select {
		case err := <-errChan:
			fmt.Print(string(colorDanger))
			c.log.Errorf("error send data to elastic %v", err)
			fmt.Print(string(colorReset))
		default:
			close(errChan)
		}
	}

	if dl.StatusCode >= 200 && dl.StatusCode < 400 {
		fmt.Println(string(colorSuccess), strings.Repeat("=", 60))
		level = logrus.InfoLevel
	} else if dl.StatusCode >= 400 && dl.StatusCode < 500 {
		fmt.Println(string(colorWarning), strings.Repeat("=", 60))
		level = logrus.WarnLevel
	} else {
		fmt.Println(string(colorDanger), strings.Repeat("=", 60))
		level = logrus.ErrorLevel
	}
	c.log.WithField("incoming_log", dl).Log(level, "apps")
}
