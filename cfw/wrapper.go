package cfw

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/logger"
)

type Wrapper struct {
	Config   *models.Config
	Resource *cloudflare.ResourceContainer

	context  context.Context
	cancelFn context.CancelFunc
	api      *cloudflare.API
	logger   *logger.Logger
}

func (w *Wrapper) ListRecords(t models.IPType) (map[string]cloudflare.DNSRecord, error) {
	var result = make(map[string]cloudflare.DNSRecord)
	records, _, err := w.api.ListDNSRecords(w.context, w.Resource, cloudflare.ListDNSRecordsParams{
		Type: string(t.RecordType),
	})
	if err != nil {
		w.cancelFn()
		return result, err
	}

	for _, record := range records {
		result[record.Name] = record
	}

	return result, err
}

func (w *Wrapper) CreateRecord(t models.IPType, query models.SubDomain) error {
	w.logger.Debug("creating record '%s'", query.Name)
	var resp, err = w.api.CreateDNSRecord(w.context, w.Resource, cloudflare.CreateDNSRecordParams{
		Name:    query.Name,
		Type:    string(t.RecordType),
		Content: query.IP,
		TTL:     GetTTL(query.TTL, w.logger.Extend("utils")),
		Proxied: &query.Proxied,
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("cannot create record %s because %v", query.Name, resp.Errors)
	}

	w.logger.Info("created record '%s'", query.Name)
	return nil
}

func (w *Wrapper) UpdateRecord(t models.IPType, query models.SubDomain) error {
	w.logger.Debug("updating record '%s'", query.Name)
	var err = w.api.UpdateDNSRecord(w.context, w.Resource, cloudflare.UpdateDNSRecordParams{
		ID:      query.Id,
		Name:    query.Name,
		Type:    string(t.RecordType),
		Content: query.IP,
		TTL:     GetTTL(query.TTL, w.logger.Extend("utils")),
		Proxied: &query.Proxied,
	})
	if err != nil {
		return err
	}

	w.logger.Info("updated record '%s'", query.Name)
	return nil
}

func (w *Wrapper) DeleteRecord(t models.IPType, query models.SubDomain) error {
	w.logger.Debug("deleting record %s", query.Id)
	var err = w.api.DeleteDNSRecord(w.context, w.Resource, query.Id)
	if err != nil {
		return err
	}

	w.logger.Info("deleted record '%s'", query.Name)
	return nil
}

func (w *Wrapper) UpsertRecords() {
	if len(w.Config.SubDomains) <= 0 {
		w.logger.Errorf("cannot found subdomains to upsert")
		return
	}

	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		w.upsertRecords(models.IPV4)
	}()
	go func() {
		defer wg.Done()

		w.upsertRecords(models.IPV6)
	}()

	wg.Wait()
}

func (w *Wrapper) upsertRecords(t models.IPType) {
	var setting = w.Config.Settings.GetIPSettings(t)
	var ip, err = GetPublicIP(*setting)
	if err != nil {
		w.logger.Warnf("cannot get public ip: %v", err)
		return
	}
	w.logger.Debug("%s: %s", t.Name, ip)
	if ip == nil {
		// Skip disable record
		w.logger.Info("disabled %s records, skipped", t.Name)
		return
	}

	if !t.Check(ip) {
		w.logger.Warnf("received ip(%s) is not %s", ip, t.Name)
		return
	}

	records, err := w.ListRecords(t)
	if err != nil {
		w.logger.Warnf("cannot get dns records: %v", err)
		return
	}
	w.logger.Debug("found %d dns records for type %s", len(records), t.RecordType)

	for _, subdomain := range w.Config.SubDomains {
		w.logger.Debug("processing '%s' record", subdomain.Name)
		var name = GetFullDomain(subdomain.Name, w.Config.Secrets.ZoneName)
		var record, isUpdate = records[name]

		var query = subdomain
		if query.Id == "" && isUpdate {
			query.Id = record.ID
		}
		if query.TTL == "" {
			query.TTL = w.Config.Settings.Ttl
		}
		if query.IP == "" {
			query.IP = ip.String()
		}
		// override name with correct full domain name
		query.Name = name

		if isUpdate && !subdomain.Enabled {
			var err = w.DeleteRecord(t, query)
			if err != nil {
				w.logger.Errorf("delete record '%s' failed because %v", query.Name, err)
			}
		} else if isUpdate && subdomain.Enabled {
			if IsChange(record, query, w.logger) {
				var err = w.UpdateRecord(t, query)
				if err != nil {
					w.logger.Errorf("update record '%s' failed because %v", query.Name, err)
				}
			} else {
				w.logger.Warnf("skipped record '%s' because same content", query.Name)
			}
		} else if !isUpdate && subdomain.Enabled {
			var err = w.CreateRecord(t, query)
			if err != nil {
				w.logger.Errorf("create record '%s' failed because %v", query.Name, err)
			}
		}
	}
}
