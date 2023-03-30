package models

import (
	"fmt"

	"github.com/kamontat/cloudflare-ddns/cloudflare"
	"github.com/kamontat/cloudflare-ddns/configs"
)

type TunnelIngress struct {
	Tunnel  *cloudflare.TunnelRecord
	Name    string
	Path    string
	Service string
}

func ToIngress(entity *Entity, tunnel *cloudflare.TunnelRecord) (*TunnelIngress, error) {
	// Only build ingress for tunnel mode entity
	if entity.ModeKey != configs.TunnelModeKey {
		return nil, nil
	}

	var ingress = &TunnelIngress{
		Tunnel:  tunnel,
		Name:    entity.Name,
		Path:    entity.Tunnel.Path,
		Service: entity.Tunnel.Service,
	}
	if ingress.Service == "" {
		return nil, fmt.Errorf("tunnel-service cannot be empty if ingress is enabled")
	}

	return ingress, nil
}

func ToTunnelConfigs(ingresses []*TunnelIngress, setting configs.Setting) (configs []*cloudflare.TunnelConfig) {
	if len(ingresses) <= 0 {
		return
	}

	var mapper = make(map[string]*cloudflare.TunnelConfig)
	for _, ingress := range ingresses {
		var config *cloudflare.TunnelConfig
		if value, ok := mapper[ingress.Tunnel.Id]; ok {
			config = value
		} else {
			config = &cloudflare.TunnelConfig{
				Record:          ingress.Tunnel,
				Ingresses:       make([]*cloudflare.TunnelConfigIngress, 0),
				CatchallService: setting.GetTunnelCatchallService(),
			}
		}

		config.Ingresses = append(config.Ingresses, &cloudflare.TunnelConfigIngress{
			Name:    ingress.Name,
			Path:    ingress.Path,
			Service: ingress.Service,
		})

		mapper[ingress.Tunnel.Id] = config
	}

	for _, value := range mapper {
		configs = append(configs, value)
	}
	return
}
