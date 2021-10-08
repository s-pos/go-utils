package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/s-pos/go-utils/logger"
	"github.com/s-pos/go-utils/utils/monitoring"
	"github.com/sirupsen/logrus"
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

	if dl.StatusCode >= 200 && dl.StatusCode < 400 {
		level = logrus.InfoLevel
	} else if dl.StatusCode >= 400 && dl.StatusCode < 500 {
		level = logrus.WarnLevel
	} else {
		level = logrus.ErrorLevel
	}

	monitoring.Prometheus().Record(dl.StatusCode, dl.RequestMethod, dl.Endpoint, dl.ResponseMessage, time.Since(dl.TimeStart))
	if elasticStatus {
		select {
		case err := <-errChan:
			c.log.Errorf("error send data to elastic %v", err)
		default:
			close(errChan)
		}
	}
	c.log.WithField("incoming_log", dl).Log(level, "apps")
}
