package configs

type ModeKey string

const (
	IPModeKey     ModeKey = "ip"
	IPV4ModeKey   ModeKey = "ipv4"
	IPV6ModeKey   ModeKey = "ipv6"
	TunnelModeKey ModeKey = "tunnel"
)

var (
	ModeKeys = []string{
		string(IPModeKey),
		string(IPV4ModeKey),
		string(IPV6ModeKey),
		string(TunnelModeKey),
	}
)

type DefaultEntity struct {
	ModeKey       *ModeKey `json:"mode-key" yaml:"mode-key"`
	TunnelName    *string  `json:"tunnel-name" yaml:"tunnel-name"`
	TunnelDomain  *string  `json:"tunnel-domain" yaml:"tunnel-domain"`
	TunnelPath    *string  `json:"tunnel-path" yaml:"tunnel-path"`
	TunnelService *string  `json:"tunnel-service" yaml:"tunnel-service"`
	Enabled       *bool    `json:"enabled" yaml:"enabled"`
	Proxied       *bool    `json:"proxied" yaml:"proxied"`
	Ttl           *string  `json:"ttl" yaml:"ttl"`
}

type Entity struct {
	Name *string `json:"name" yaml:"name"`
	DefaultEntity
}
