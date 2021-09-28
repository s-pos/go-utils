package request

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

// Methods interface for make request with all http.Method
type Methods interface {
	Post(ctx context.Context, payload []byte) ([]byte, int, error)
	Get(ctx context.Context) ([]byte, int, error)
	Put(ctx context.Context, payload []byte) ([]byte, int, error)
	Delete(ctx context.Context, payload []byte) ([]byte, int, error)
	WithTimeout(d int64)
	WithBasicAuth(username, password string)
}

func buf(p []byte) io.ReadCloser {
	if p != nil {
		r := bytes.NewReader(p)
		return ioutil.NopCloser(r)
	}
	return nil
}

func (r *component) do(payload []byte, method string) ([]byte, int, error) {

	req, err := http.NewRequest(method, r.url, buf(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Set header
	if r.header != nil {
		for k, v := range r.header {
			req.Header.Set(k, v)
		}
	}

	// Set basic auth if exist
	if r.basicAuth.set {
		req.SetBasicAuth(r.basicAuth.username, r.basicAuth.password)
	}

	// Set timeout if exist
	if r.timeout.set {
		cons, cancel := context.WithTimeout(context.Background(), r.timeout.duration)
		req = req.WithContext(cons)
		defer cancel()
	}

	// Do request
	res, err := r.client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	// Read result
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return result, res.StatusCode, nil
}

// Request for create a new data request
func (r *component) Request(header map[string]string, url string) Methods {
	r.header = header
	r.url = url
	return r
}

// Post is request with method POST
func (r *component) Post(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodPost)
}

// Get is request with method GET
func (r *component) Get(ctx context.Context) ([]byte, int, error) {
	return r.wrapper(ctx, nil, http.MethodGet)
}

// Delete is request with method DELETE
func (r *component) Delete(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodDelete)
}

// Put is request with method PUT
func (r *component) Put(ctx context.Context, payload []byte) ([]byte, int, error) {
	return r.wrapper(ctx, payload, http.MethodPut)
}
