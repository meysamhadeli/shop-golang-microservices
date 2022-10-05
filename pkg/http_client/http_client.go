package http_client

import (
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
	idleConnTimeout       = 120 * time.Second
	responseHeaderTimeout = 5 * time.Second
)

func NewHttpClient() *http.Client {

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

	client := &http.Client{
		Transport: otelhttp.NewTransport(transport), // use custom transport open-telemetry for tracing http-client
		Timeout:   timeout,
	}

	return client
}
