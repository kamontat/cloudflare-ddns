package clients

import (
	"crypto/tls"
	"net/http"
	"time"
)

func New() *http.Client {
	var transport *http.Transport = http.DefaultTransport.(*http.Transport).Clone()

	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
	}

	var client = &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	return client
}

var Default = New()
