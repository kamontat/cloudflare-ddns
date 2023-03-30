package cloudflare

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/utils"
	"github.com/kc-workspace/go-lib/caches"
)

type Cloudflare struct {
	ZoneName          string
	ZoneIdentifier    *cloudflare.ResourceContainer
	AccountName       string
	AccountIdentifier *cloudflare.ResourceContainer
	api               *cloudflare.API
	context           context.Context
	cancelFn          context.CancelFunc
	cache             *caches.Service
}

func (c *Cloudflare) ListDNSRecords() (
	records map[string]map[DNSRecordType]*DNSRecord,
	err error,
) {
	records = make(map[string]map[DNSRecordType]*DNSRecord)
	rawRecords, _, err := c.api.ListDNSRecords(
		c.context,
		c.ZoneIdentifier,
		cloudflare.ListDNSRecordsParams{},
	)
	if err != nil {
		return
	}

	for _, rec := range rawRecords {
		// Filter only support dns record types
		if utils.ContainValue(DNSRecordTypes, rec.Type) {
			var record = ToDNSRecord(rec, c.ZoneName)
			var nested = make(map[DNSRecordType]*DNSRecord)
			if value, ok := records[record.Name]; ok {
				nested = value
			}
			nested[record.Type] = &record
			records[record.Name] = nested
		}
	}

	return
}

func (c *Cloudflare) ListTunnelRecords() (records []TunnelRecord, err error) {
	tunnels, _, err := c.api.ListTunnels(
		c.context,
		c.AccountIdentifier,
		cloudflare.TunnelListParams{
			// https://github.com/cloudflare/cloudflare-go/issues/1247
			ResultInfo: cloudflare.ResultInfo{
				PerPage: 1000,
			},
		},
	)
	if err != nil {
		return
	}

	for _, tunnel := range tunnels {
		records = append(records, ToTunnelRecord(tunnel))
	}
	return
}

func (c *Cloudflare) GetTunnelRecord(name string) (*TunnelRecord, error) {
	var cacheKey = fmt.Sprintf("cloudflare.tunnel.record.%s", name)
	if c.cache.Has(cacheKey) {
		var cacheData = c.cache.Get(cacheKey)
		var raw, err = cacheData.Get()
		return raw.(*TunnelRecord), err
	}

	var tunnels, _, err = c.api.ListTunnels(
		c.context,
		c.AccountIdentifier,
		cloudflare.TunnelListParams{
			Name: name,
			// https://github.com/cloudflare/cloudflare-go/issues/1247
			ResultInfo: cloudflare.ResultInfo{
				PerPage: 5,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(tunnels) < 1 {
		return nil, fmt.Errorf("cannot found tunnel name %s", name)
	}

	var record = ToTunnelRecord(tunnels[0])
	c.cache.Set(cacheKey, &record, "15m")
	return &record, nil
}

func (c *Cloudflare) UpdateTunnelConfig(config *TunnelConfig) (err error) {
	var ingresses = make([]cloudflare.UnvalidatedIngressRule, 0)
	for _, ingress := range config.Ingresses {
		ingresses = append(ingresses, cloudflare.UnvalidatedIngressRule{
			Hostname: utils.BuildRecordName(ingress.Name, c.ZoneName),
			Path:     ingress.Path,
			Service:  ingress.Service,
		})
	}

	// Add catch-all
	ingresses = append(ingresses, cloudflare.UnvalidatedIngressRule{
		Service: config.CatchallService,
	})

	var params = cloudflare.TunnelConfigurationParams{
		TunnelID: config.Record.Id,
		Config: cloudflare.TunnelConfiguration{
			Ingress: ingresses,
		},
	}

	_, err = c.api.UpdateTunnelConfiguration(
		c.context,
		c.AccountIdentifier,
		params,
	)
	return
}

func (c *Cloudflare) UpsertDNSRecord(record *DNSRecord) error {
	if record.Id == "" {
		return c.InsertDNSRecord(record)
	} else {
		return c.UpdateDNSRecord(record)
	}
}

func (c *Cloudflare) InsertDNSRecord(record *DNSRecord) error {
	if record.Id != "" {
		return fmt.Errorf("your dns record has been created, use update instead")
	}
	var now = time.Now().UTC().Format(time.RFC3339)
	var comment = fmt.Sprintf("insert by cfddns at %s", now)
	var _, err = c.api.CreateDNSRecord(
		c.context,
		c.ZoneIdentifier,
		cloudflare.CreateDNSRecordParams{
			Name:    utils.BuildRecordName(record.Name, c.ZoneName),
			Content: record.Content,
			Type:    string(record.Type),
			TTL:     record.GetTTL(),
			Proxied: &record.Proxied,
			Comment: comment,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloudflare) UpdateDNSRecord(record *DNSRecord) error {
	if record.Id == "" {
		return fmt.Errorf("your dns record never created before, use insert instead")
	}
	var now = time.Now().UTC().Format(time.RFC3339)
	var comment = fmt.Sprintf("update by cfddns at %s", now)
	var err = c.api.UpdateDNSRecord(
		c.context,
		c.ZoneIdentifier,
		cloudflare.UpdateDNSRecordParams{
			ID:      record.Id,
			Name:    utils.BuildRecordName(record.Name, c.ZoneName),
			Type:    string(record.Type),
			Content: record.Content,
			TTL:     record.GetTTL(),
			Proxied: &record.Proxied,
			Comment: comment,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloudflare) DeleteDNSRecord(id string) error {
	if id == "" {
		return fmt.Errorf("cannot found dns record id empty")
	}

	var err = c.api.DeleteDNSRecord(
		c.context,
		c.ZoneIdentifier,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloudflare) Close() {
	c.cancelFn()
}
