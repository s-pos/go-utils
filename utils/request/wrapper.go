package request

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/s-pos/go-utils/logger"
	"github.com/s-pos/go-utils/utils/monitoring"
)

// wrapper is doing request to 3rd Party
func (r *component) wrapper(ctx context.Context, payload []byte, method string) ([]byte, int, error) {
	start := time.Now()
	var tp logger.ThirdParty

	tp.Method = method
	tp.URL = r.url
	tp.RequestHeader = headerToString(&r.header)

	if payload != nil {
		tp.RequestBody = string(payload)
	}

	data, status, err := r.do(payload, method)

	tp.StatusCode = status
	if err != nil {
		tp.Response = err.Error()
	}

	if data != nil {
		var resBody interface{}

		err = json.Unmarshal(data, &resBody)
		if err != nil {
			tp.Response = "response is empty"
		} else {
			tp.Response = resBody
		}
	}

	since := time.Since(start)
	tp.ExecTime = since.Seconds()

	/* storing third party request and response to context and prometheus */
	tp.Store(ctx)
	monitoring.Prometheus().Record(tp.StatusCode, tp.Method, tp.URL, "outgoing_request", since)

	return data, status, err
}

func headerToString(h *map[string]string) string {
	var header string

	if h == nil {
		return header
	}

	for key, v := range *h {
		header += fmt.Sprintf(` %v: %v `, key, v)
	}
	return header
}
