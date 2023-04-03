package cloudflare

import (
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func get[T any](i T) *T {
	return &i
}

func buildOriginConfig() cloudflare.OriginRequestConfig {
	// All duration are fixed to second unit
	// meaning time.Duration(15) is 15 second
	return cloudflare.OriginRequestConfig{
		ConnectTimeout:       get(time.Duration(15)),
		TLSTimeout:           get(time.Duration(15)),
		TCPKeepAlive:         get(time.Duration(15)),
		KeepAliveConnections: get(50),
		KeepAliveTimeout:     get(time.Duration(120)),
	}
}
