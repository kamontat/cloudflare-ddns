package cloudflare

import "github.com/cloudflare/cloudflare-go"

type TunnelConfigIngress struct {
	Name    string
	Path    string
	Service string
}

type TunnelConfig struct {
	Record          *TunnelRecord
	Ingresses       []*TunnelConfigIngress
	CatchallService string
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
