package cloudflare

import "github.com/cloudflare/cloudflare-go"

type TunnelConfigIngress struct {
	Host    string
	Path    string
	Service string
}

type TunnelConfig struct {
	Record          TunnelRecord
	Ingresses       []TunnelConfigIngress
	CatchallService string
}

func ToTunnelConfig(record TunnelRecord, config cloudflare.TunnelConfiguration) TunnelConfig {
	var catchall string
	var ingresses = make([]TunnelConfigIngress, 0)
	for _, ingress := range config.Ingress {
		if ingress.Hostname == "" && ingress.Path == "" && ingress.Service != "" {
			catchall = ingress.Service
		} else {
			ingresses = append(ingresses, TunnelConfigIngress{
				Host:    ingress.Hostname,
				Path:    ingress.Path,
				Service: ingress.Service,
			})
		}
	}

	return TunnelConfig{
		Record:          record,
		Ingresses:       ingresses,
		CatchallService: catchall,
	}
}
