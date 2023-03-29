package models

import (
	"fmt"
	"time"

	"github.com/kamontat/cloudflare-ddns/cloudflare"
	"github.com/kamontat/cloudflare-ddns/configs"
	"github.com/kamontat/cloudflare-ddns/utils"
)

type Tunnel struct {
	TunnelName   string
	TunnelDomain string
}

type Entity struct {
	Name    string
	ModeKey configs.ModeKey
	Enabled bool
	Proxied bool
	TTL     int
	Query   *configs.Query
	Tunnel
}

func (e *Entity) GetType() cloudflare.DNSRecordType {
	switch e.ModeKey {
	case configs.IPV4ModeKey:
		return cloudflare.A
	case configs.IPV6ModeKey:
		return cloudflare.AAAA
	case configs.TunnelModeKey:
		return cloudflare.CNAME
	default:
		return cloudflare.TXT
	}
}

func (e *Entity) GetContent(cf *cloudflare.Cloudflare) (string, error) {
	switch e.ModeKey {
	case configs.IPV4ModeKey, configs.IPV6ModeKey:
		return e.Query.Query()
	case configs.TunnelModeKey:
		tunnel, err := cf.GetTunnelRecord(e.TunnelName)
		if err != nil {
			return "", err
		}

		return tunnel.GetURL(e.TunnelDomain), nil
	default:
		return "", fmt.Errorf("cannot found mode-key to resolve content")
	}
}

func (e *Entity) GetTTL() int {
	return utils.BuildTTL(e.TTL)
}

// Compare entity with record.
func (e *Entity) Compare(record *cloudflare.DNSRecord, cf *cloudflare.Cloudflare) (*cloudflare.DNSRecord, error) {
	if !e.Enabled {
		return nil, nil
	}
	var rec = cloudflare.DNSRecord{}
	if record != nil {
		rec = *record
	}

	var isChanged = false
	if rec.Name != e.Name {
		rec.Name = e.Name
		isChanged = true
	}
	if rec.Type != e.GetType() {
		rec.Type = e.GetType()
		isChanged = true
	}
	content, err := e.GetContent(cf)
	if err != nil {
		return nil, err
	}
	if rec.Content != content {
		rec.Content = content
		isChanged = true
	}
	if rec.Proxied != e.Proxied {
		rec.Proxied = e.Proxied
		isChanged = true
	}
	if rec.GetTTL() != e.GetTTL() {
		rec.Ttl = e.GetTTL()
		isChanged = true
	}

	if !isChanged {
		return nil, nil
	}

	return &rec, nil
}

func ToEntities(config configs.Entity, def *configs.DefaultEntity, setting configs.Setting) (entities []*Entity, err error) {
	var base = utils.GetOrElse(def, configs.DefaultEntity{})

	modeKey, _ := utils.GetOr("entity.mode-key", config.ModeKey, base.ModeKey, configs.IPModeKey)
	if !utils.ContainValue(configs.ModeKeys, string(modeKey)) {
		err = fmt.Errorf("invalid mode-key: %s", modeKey)
		return
	}

	var tunnel = &Tunnel{}
	if modeKey == configs.TunnelModeKey {
		var name string
		name, err = utils.GetOr("entity.tunnel-name", config.TunnelName, base.TunnelName, "")
		if err != nil {
			return
		}
		if name == "" {
			err = fmt.Errorf("tunnel name must be defined if your mode-key is tunnel")
			return
		}
		tunnel.TunnelName = name

		var domain string
		domain, err = utils.GetOr("entity.tunnel-domain", config.TunnelDomain, base.TunnelDomain, "")
		if err != nil {
			return
		}
		tunnel.TunnelDomain = domain
	}

	var query *configs.Query = nil
	if modeKey == configs.IPV4ModeKey {
		query = setting.Ipv4
	} else if modeKey == configs.IPV6ModeKey {
		query = setting.Ipv6
	}

	name, err := utils.GetOr("entity.name", config.Name, nil, "")
	if err != nil {
		return
	}

	enabled, _ := utils.GetOr("entity.enabled", config.Enabled, base.Enabled, false)
	proxied, _ := utils.GetOr("entity.proxied", config.Proxied, base.Proxied, false)
	ttl, _ := utils.GetOr("entity.ttl", config.Ttl, base.Ttl, "0")
	dur, err := time.ParseDuration(ttl)
	if err != nil {
		err = fmt.Errorf("invalid ttl duration: %s", err.Error())
		return
	}

	// separate ip mode-key to ipv4 and ipv6
	if modeKey == configs.IPModeKey {
		entities = []*Entity{{
			Name:    name,
			ModeKey: configs.IPV4ModeKey,
			Enabled: enabled,
			Proxied: proxied,
			TTL:     int(dur.Seconds()),
			Query:   setting.Ipv4,
			Tunnel:  *tunnel,
		}, {
			Name:    name,
			ModeKey: configs.IPV6ModeKey,
			Enabled: enabled,
			Proxied: proxied,
			TTL:     int(dur.Seconds()),
			Query:   setting.Ipv6,
			Tunnel:  *tunnel,
		}}
	} else {
		entities = []*Entity{{
			Name:    name,
			ModeKey: modeKey,
			Enabled: enabled,
			Proxied: proxied,
			TTL:     int(dur.Seconds()),
			Query:   query,
			Tunnel:  *tunnel,
		}}
	}

	return
}
