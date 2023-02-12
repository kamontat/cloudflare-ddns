package models

type SubDomain struct {
	Enabled bool
	Id      string
	Name    string
	IP      string
	TTL     string
	Proxied bool
}
