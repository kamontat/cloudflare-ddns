package models

type IPQuerySettings struct {
	Url       string
	Format    string
	Separator string
	Key       string
}

type IPSettings struct {
	Enabled bool
	Query   IPQuerySettings
}

type Settings struct {
	Ipv4 IPSettings
	Ipv6 IPSettings
	// Ttl is time.Duration format
	Ttl   string
	Purge bool
}
