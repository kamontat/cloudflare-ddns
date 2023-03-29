package configs

type Setting struct {
	Defaults *DefaultEntity `json:"defaults" yaml:"defaults"`
	Ipv4     *Query         `json:"ipv4" yaml:"ipv4"`
	Ipv6     *Query         `json:"ipv6" yaml:"ipv6"`
}
