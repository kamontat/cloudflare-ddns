package core

import (
	"github.com/kamontat/cloudflare-ddns/cloudflare"
	"github.com/kamontat/cloudflare-ddns/configs"
	"github.com/kamontat/cloudflare-ddns/models"
	cmodels "github.com/kc-workspace/go-lib/commandline/models"
	"github.com/kc-workspace/go-lib/logger"
)

type Core struct {
	meta       *cmodels.Metadata
	config     *configs.Config
	cloudflare *cloudflare.Cloudflare
	logger     *logger.Logger
}

func (c *Core) Start() error {
	records, err := c.cloudflare.ListDNSRecords()
	if err != nil {
		return err
	}

	var ingresses = make([]*models.TunnelIngress, 0)

	c.logger.Debug("listed dns-records size: %d", len(records))
	for _, raw := range c.config.Entities {
		var entities, err = models.ToEntities(
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
			if !entity.Enabled {
				c.logger.Info("disabled: %s entity (%s)", entity.Name, entity.ModeKey)
				continue
			}

			// Update DNS Records
			var rec = records[entity.Name][entity.GetType()]
			c.logger.Debug("getting original record: %#v", rec)
			var record, err = entity.Compare(rec, c.cloudflare)
			c.logger.Debug("built modified record: %#v", record)
			if err != nil {
				c.logger.Error(err)
			} else {
				// DNS Record
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
					}
				} else {
					c.logger.Info("do nothing on dns-record: %s", entity.Name)
				}
			}

			// Tunnel ingress
			if entity.Tunnel != nil && c.config.Settings.GetTunnelIngress() {
				tunnel, err := c.cloudflare.GetTunnelRecord(entity.Tunnel.Name)
				if err != nil {
					c.logger.Error(err)
					continue
				}

				ingress, err := models.ToIngress(entity, tunnel)
				if err != nil {
					c.logger.Error(err)
					continue
				}

				if ingress != nil {
					ingresses = append(ingresses, ingress)
				}
			}
		}
	}

	if c.config.Settings.GetTunnelIngress() {
		var tunnelConfigs = models.ToTunnelConfigs(ingresses, c.config.Settings)
		for _, config := range tunnelConfigs {
			c.logger.Info("updating tunnel: %s (size=%d)", config.Record.Name, len(config.Ingresses))
			if err = c.cloudflare.UpdateTunnelConfig(config); err != nil {
				c.logger.Error(err)
			}
		}
	}
	return nil
}

func (c *Core) Close() {
	c.cloudflare.Close()
}
