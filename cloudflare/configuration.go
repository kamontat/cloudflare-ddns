package cloudflare

import (
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func get[T any](i T) *T {
	return &i
}

func buildOriginConfig(zone string) cloudflare.OriginRequestConfig {
	return cloudflare.OriginRequestConfig{
		ConnectTimeout:       get(10 * time.Second),
		TLSTimeout:           get(10 * time.Second),
		TCPKeepAlive:         get(30 * time.Second),
		KeepAliveConnections: get(30),
		KeepAliveTimeout:     get(2 * time.Minute),
		OriginServerName:     get("*." + zone),
		NoTLSVerify:          get(true),
	}
}
