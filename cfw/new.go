package cfw

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/logger"
)

func New(config *models.Config, log *logger.Logger) (w *Wrapper, err error) {
	var ctx, cancelFn = context.WithCancel(context.Background())

	var dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	var client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          5,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	api, err := cloudflare.NewWithAPIToken(
		config.Secrets.ApiToken,
		cloudflare.HTTPClient(client),
	)

	if err != nil {
		cancelFn()
		return
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(config.Secrets.ZoneName)
	if err != nil {
		cancelFn()
		return
	}

	w = &Wrapper{
		Config:   config,
		Resource: cloudflare.ZoneIdentifier(id),
		context:  ctx,
		cancelFn: cancelFn,
		api:      api,
		logger:   log.Extend("cloudflare"),
	}

	return
}
