package core

import (
	"github.com/kamontat/cloudflare-ddns/cloudflare"
	"github.com/kamontat/cloudflare-ddns/configs"
	dmodels "github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/commandline/models"
	"github.com/kc-workspace/go-lib/logger"
)

type Core struct {
	meta       *models.Metadata
	config     *configs.Config
	cloudflare *cloudflare.Cloudflare
	logger     *logger.Logger
}

func (c *Core) Start() error {
	records, err := c.cloudflare.ListDNSRecords()
	if err != nil {
		return err
	}

	c.logger.Debug("listed dns-records size: %d", len(records))
	for _, raw := range c.config.Entities {
		var entities, err = dmodels.ToEntities(
			raw,
			c.config.Settings.Defaults,
			c.config.Settings,
		)
		if err != nil {
			c.logger.Error(err)
			continue
		}

		for _, entity := range entities {
			c.logger.Debug("processing entity: %#v", entity)

			var rec = records[entity.Name][entity.GetType()]
			c.logger.Debug("getting original record: %#v", rec)
			var record, err = entity.Compare(rec, c.cloudflare)
			c.logger.Debug("built modified record: %#v", record)
			if err != nil {
				c.logger.Error(err)
				continue
			}

			if !entity.Enabled && rec != nil {
				c.logger.Info("deleting dns-record: %s", rec.Id)
				err = c.cloudflare.DeleteDNSRecord(rec.Id)
				if err != nil {
					c.logger.Error(err)
				}
			} else if record != nil {
				c.logger.Info("upserting dns-record: %s (%s)", record.Id, record.Name)
				err = c.cloudflare.UpsertDNSRecord(record)
				if err != nil {
					c.logger.Error(err)
					continue
				}
			} else {
				c.logger.Info("do nothing: %s", entity.Name)
			}
		}
	}

	return nil
}
