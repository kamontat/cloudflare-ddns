package cloudflare

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type TunnelStatus string

const (
	HEALTHY  = "healthy"
	DOWN     = "down"
	INACTIVE = "active"
)

type TunnelRecord struct {
	Id     string
	Name   string
	Status TunnelStatus
}

func (r *TunnelRecord) GetURL(domain string) string {
	return fmt.Sprintf("%s.%s", r.Id, domain)
}

func ToTunnelRecord(tunnel cloudflare.Tunnel) TunnelRecord {
	return TunnelRecord{
		Id:     tunnel.ID,
		Name:   tunnel.Name,
		Status: TunnelStatus(tunnel.Status),
	}
}
