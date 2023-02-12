package models

import "net"

type RecordType string

const (
	A    RecordType = "A"
	AAAA RecordType = "AAAA"
)

type IPType struct {
	Name       string
	RecordType RecordType
	Check      func(ip *net.IP) bool
}

var (
	IPV4 IPType = IPType{
		Name:       "ipv4",
		RecordType: A,
		Check: func(ip *net.IP) bool {
			return ip.To4() != nil
		},
	}

	IPV6 IPType = IPType{
		Name:       "ipv6",
		RecordType: AAAA,
		Check: func(ip *net.IP) bool {
			return ip.To4() == nil
		},
	}
)
