package main

import "time"

type SubDomain struct {
	Name    string
	Enabled bool
	Proxied bool
}

type Config struct {
	Ipv4       bool
	Ipv6       bool
	Ttl        time.Duration
	Purge      bool
	SubDomains []SubDomain
}
