package request

import (
	"net/http"
	"time"
)

// component data request
type component struct {
	url    string
	header map[string]string
	timeout
	client    *http.Client
	basicAuth basicAuth
}

type timeout struct {
	duration time.Duration
	set      bool
}

type basicAuth struct {
	username string
	password string
	set      bool
}

// WithTimeout will set the timeout of request outgoing
func (r *component) WithTimeout(d int64) {
	r.timeout.duration = time.Duration(d) * time.Second
	r.timeout.set = true

}

// WithBasicAuth is setting basic auth
func (r *component) WithBasicAuth(username, password string) {
	r.basicAuth.username = username
	r.basicAuth.password = password
	r.basicAuth.set = true
}

/*Client is client to make new clients */
type Client interface {
	Request(header map[string]string, url string) Methods
}

/*NewHTTPClient request with new http client*/
func NewHTTPClient(client *http.Client) Client {
	if client == nil {
		client = http.DefaultClient
	}
	c := new(component)
	c.client = client
	return c
}
