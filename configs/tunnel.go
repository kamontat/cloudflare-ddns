package configs

type Tunnel struct {
	Ingress         *bool   `json:"ingress" yaml:"ingress"`
	CatchallService *string `json:"catchall-service" yaml:"catchall-service"`
}

var (
	DefaultCatchallService = "http_status:404"
)
