package cfw

import (
	"crypto/tls"
	"net/http"
)

func NewClient() *http.Client {
	var transport *http.Transport = http.DefaultTransport.(*http.Transport).Clone()

	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
	}

	var client = &http.Client{
		Transport: transport,
	}

	return client
}

var DefaultClient = NewClient()
