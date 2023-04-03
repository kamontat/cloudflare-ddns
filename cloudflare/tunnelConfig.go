package cloudflare

import "github.com/cloudflare/cloudflare-go"

type TunnelConfigIngress struct {
	Name    string
	Path    string
	Service string
}

func (c *TunnelConfigIngress) IsSame(other *TunnelConfigIngress) bool {
	return c.Name == other.Name &&
		c.Path == other.Path &&
		c.Service == other.Service
}

type TunnelConfig struct {
	Record          *TunnelRecord
	Ingresses       []*TunnelConfigIngress
	CatchallService string
}

func (c *TunnelConfig) IsSame(other *TunnelConfig) bool {
	if c.Record.Id != other.Record.Id ||
		c.Record.Name != other.Record.Name {
		return false
	}

	if len(c.Ingresses) != len(other.Ingresses) {
		return false
	}

	for i := 0; i < len(c.Ingresses); i++ {
		var ci = c.Ingresses[i]
		var oi = other.Ingresses[i]
		if !ci.IsSame(oi) {
			return false
		}
	}

	return c.CatchallService == other.CatchallService
}

// For support origin request on ingress rule level
// ref: https://github.com/cloudflare/cloudflare-go/pull/1138
type TunnelConfigurationIngress struct {
	Hostname      string                          `json:"hostname,omitempty"`
	Path          string                          `json:"path,omitempty"`
	Service       string                          `json:"service,omitempty"`
	OriginRequest *cloudflare.OriginRequestConfig `json:"originRequest,omitempty"`
}

// ref: https://github.com/cloudflare/cloudflare-go/pull/1138
type TunnelConfiguration struct {
	Ingress       []TunnelConfigurationIngress   `json:"ingress,omitempty"`
	WarpRouting   *cloudflare.WarpRoutingConfig  `json:"warp-routing,omitempty"`
	OriginRequest cloudflare.OriginRequestConfig `json:"originRequest,omitempty"`
}

// ref: https://github.com/cloudflare/cloudflare-go/pull/1138
type TunnelConfigurationParams struct {
	TunnelID string              `json:"-"`
	Config   TunnelConfiguration `json:"config,omitempty"`
}
