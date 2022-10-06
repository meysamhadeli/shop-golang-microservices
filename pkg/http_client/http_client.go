package http_client

import (
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"time"
)

const (
	timeout               = 5 * time.Second
	dialContextTimeout    = 5 * time.Second
	tLSHandshakeTimeout   = 5 * time.Second
	xaxIdleConns          = 20
	maxConnsPerHost       = 40
	retryCount            = 3
	retryWaitTime         = 300 * time.Millisecond
	idleConnTimeout       = 120 * time.Second
	responseHeaderTimeout = 5 * time.Second
)

func NewHttpClient() *resty.Client {

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: dialContextTimeout,
		}).DialContext,
		TLSHandshakeTimeout:   tLSHandshakeTimeout,
		MaxIdleConns:          xaxIdleConns,
		MaxConnsPerHost:       maxConnsPerHost,
		IdleConnTimeout:       idleConnTimeout,
		ResponseHeaderTimeout: responseHeaderTimeout,
	}

	client := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetTransport(otelhttp.NewTransport(transport)) // use custom transport open-telemetry for tracing http-client

	return client
}
