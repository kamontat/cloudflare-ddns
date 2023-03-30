package cloudflare

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
