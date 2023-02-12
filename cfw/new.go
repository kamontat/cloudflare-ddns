package cfw

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/logger"
)

func New(config *models.Config, log *logger.Logger) (w *Wrapper, err error) {
	var ctx, cancelFn = context.WithCancel(context.Background())

	api, err := cloudflare.NewWithAPIToken(config.Secrets.ApiToken)
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
