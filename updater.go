package main

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/commandline/commands"
)

func updateRecords(config *models.Config, p *commands.ExecutorParameter) (err error) {

	api, err := cloudflare.NewWithAPIToken(config.Secrets.ApiToken)
	if err != nil {
		return
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch the zone ID
	id, err := api.ZoneIDByName(config.Secrets.ZoneName)
	if err != nil {
		return
	}

	// Fetch DNS records
	records, _, err := api.ListDNSRecords(
		ctx,
		cloudflare.ZoneIdentifier(id),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
			ResultInfo: cloudflare.ResultInfo{
				Page:    1,
				PerPage: 100,
			},
		},
	)
	if err != nil {
		return
	}

	// Print user details
	p.Logger.Log("%#v", records)

	return
}
