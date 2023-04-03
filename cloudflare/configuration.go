package cloudflare

import (
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func get[T any](i T) *T {
	return &i
}

func buildOriginConfig() cloudflare.OriginRequestConfig {
	return cloudflare.OriginRequestConfig{
		ConnectTimeout:       get(15 * time.Second),
		TLSTimeout:           get(15 * time.Second),
		TCPKeepAlive:         get(15 * time.Second),
		KeepAliveConnections: get(50),
		KeepAliveTimeout:     get(2 * time.Minute),
	}
}
