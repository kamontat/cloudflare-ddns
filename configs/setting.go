package configs

type Setting struct {
	Defaults *DefaultEntity `json:"defaults" yaml:"defaults"`
	Ipv4     *Query         `json:"ipv4" yaml:"ipv4"`
	Ipv6     *Query         `json:"ipv6" yaml:"ipv6"`
	Tunnel   *Tunnel        `json:"tunnel" yaml:"tunnel"`
}

func (s *Setting) GetTunnelCatchallService() string {
	if s.Tunnel == nil {
		return DefaultCatchallService
	}
	if s.Tunnel.CatchallService == nil {
		return DefaultCatchallService
	}
	if *s.Tunnel.CatchallService == "" {
		return DefaultCatchallService
	}
	return *s.Tunnel.CatchallService
}

func (s *Setting) GetTunnelIngress() bool {
	if s.Tunnel == nil {
		return false
	}
	if s.Tunnel.Ingress == nil {
		return false
	}
	return *s.Tunnel.Ingress
}
