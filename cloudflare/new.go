package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/clients"
	"github.com/kamontat/cloudflare-ddns/configs"
)

func New(secret configs.Secret) (*Cloudflare, error) {
	var context, cancelFn = context.WithCancel(context.Background())
	api, err := cloudflare.NewWithAPIToken(
		secret.GetApiToken(),
		cloudflare.HTTPClient(clients.Default),
	)
	if err != nil {
		cancelFn()
		return nil, err
	}

	var accountName = ""
	var accountId = ""
	accounts, _, err := api.Accounts(context, cloudflare.AccountsListParams{})
	if secret.GetAccountName() == "" {
		accountId = accounts[0].ID
		accountName = accounts[0].Name
	} else {
		for _, acc := range accounts {
			if acc.ID == secret.GetAccountName() ||
				acc.Name == secret.GetAccountName() {
				accountId = acc.ID
				accountName = acc.Name
				break
			}
		}
	}
	if err != nil {
		cancelFn()
		return nil, err
	}

	id, err := api.ZoneIDByName(secret.GetZoneName())
	if err != nil {
		cancelFn()
		return nil, err
	}

	return &Cloudflare{
		ZoneName:          secret.GetZoneName(),
		ZoneIdentifier:    cloudflare.ZoneIdentifier(id),
		AccountName:       accountName,
		AccountIdentifier: cloudflare.AccountIdentifier(accountId),
		api:               api,
		context:           context,
		cancelFn:          cancelFn,
	}, nil
}
