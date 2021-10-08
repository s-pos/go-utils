package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/s-pos/go-utils/config"
)

// Initialize for init first context
func Initialize(req *http.Request, start time.Time) (*http.Request, DataLogger) {
	var (
		lock = new(Locker)
		dl   DataLogger
	)

	dl.RequestID = uuid.New().String()
	dl.Service = config.GetServiceName()
	dl.Host = req.Host
	dl.Endpoint = req.URL.Path
	dl.RequestMethod = req.Method
	dl.TimeStart = start
	dl.RequestHeader = dumpRequest(req)
	dl.RequestBody = dumpBody(req)

	ctx := context.WithValue(req.Context(), logKey, lock)

	return req.WithContext(ctx), dl
}

// dumpHeader for getting all request from Header
func dumpRequest(req *http.Request) map[string]interface{} {
	var reqHeader = make(map[string]interface{})

	for key, value := range req.Header {
		reqHeader[key] = strings.Join(value, ", ")
	}

	return reqHeader
}

// dumpBody for getting all request from payload body
func dumpBody(req *http.Request) map[string]interface{} {
	var reqBody = make(map[string]interface{})
	// exctract all payload
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return reqBody
	}

	// put again body to payload request
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	if len(buf) > 1000 {
		reqBody = map[string]interface{}{"body": string(buf[:1000])}
	} else {
		json.Unmarshal(buf, &reqBody)
	}

	return reqBody
}

// Response is record response
func Response(ctx context.Context, status int, res interface{}, errLocation interface{}, err error) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.Set(_StatusCode, status)
	value.Set(_Response, res)

	if err != nil {
		value.Set(_ErrorMessage, err.Error())
		value.Set(_ErrorLocation, errLocation.(string))
	}
}

// ResponseMessage record any description message response (e.g product succesfully fetch)
func ResponseMessage(ctx context.Context, message interface{}) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	value.Set(_ResponseMessage, message)
}
