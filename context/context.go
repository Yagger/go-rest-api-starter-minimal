package context

import (
	"net/http"
	"time"
)

// Context struct
type Context struct {
	HTTP *http.Client
	// add databases here
}

// Init method
func (c *Context) Init() {
	c.HTTP = CreateHTTPClient(10, 20)
}

// CreateHTTPClient func
func CreateHTTPClient(timeout int, maxIdleConn int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConn,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	return client
}
