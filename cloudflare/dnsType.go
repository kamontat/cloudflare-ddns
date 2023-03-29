package cloudflare

type DNSRecordType string

const (
	A     DNSRecordType = "A"
	AAAA  DNSRecordType = "AAAA"
	CNAME DNSRecordType = "CNAME"
	TXT   DNSRecordType = "TXT"
)

var (
	DNSRecordTypes = []string{
		string(A),
		string(AAAA),
		string(CNAME),
	}
)
