package utils

import (
	"strings"

	"github.com/kc-workspace/go-lib/utils"
)

// build entity name from dns record name
func BuildEntityName(name, zone string) string {
	if name == zone {
		return "@"
	}
	return strings.ReplaceAll(name, "."+zone, "")
}

// build dns record name from entity name
func BuildRecordName(name, zone string) string {
	if name == "" || name == "@" || name == "." {
		return zone
	}

	return utils.JoinString(".", name, zone)
}
