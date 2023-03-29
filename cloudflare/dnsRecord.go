package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-ddns/utils"
)

type DNSRecord struct {
	Id      string
	Name    string
	Type    DNSRecordType
	Content string
	Proxied bool
	Ttl     int
}

func (r *DNSRecord) GetTTL() int {
	return utils.BuildTTL(r.Ttl)
}

func ToDNSRecord(record cloudflare.DNSRecord, zone string) DNSRecord {
	return DNSRecord{
		Id:      record.ID,
		Name:    utils.BuildEntityName(record.Name, zone),
		Type:    DNSRecordType(record.Type),
		Content: record.Content,
		Proxied: *record.Proxied,
		Ttl:     record.TTL,
	}
}
