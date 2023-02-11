package cfw

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/configs"
	"github.com/kc-workspace/go-lib/mapper"
)

func query(query models.IPQuerySettings) (result string, err error) {
	if query.Url == "" {
		err = errors.New("query.url is required")
		return
	}
	if query.Format == "" {
		err = errors.New("query.format is required")
		return
	}

	req, err := http.Get(query.Url)
	if err != nil {
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	return convert(query.Format, body, query.Key, query.Separator)
}

func convert(format string, body []byte, key, separator string) (result string, err error) {
	switch format {
	case "text":
		result = string(body)
		return
	case "json":
		var m mapper.Mapper
		m, err = mapper.FromJson(body)
		if err != nil {
			return
		}

		result, err = m.Se(key)
		return
	case "kv":
		var lines = strings.Split(string(body), "\n")
		for _, line := range lines {
			var k, v, ok = configs.ParseOverride(strings.Replace(line, separator, "=", 1))
			if ok && key == k {
				result = v
				return
			}
		}

		err = fmt.Errorf("cannot found key '%s' in response", key)
		return
	default:
		err = fmt.Errorf("invalid format name '%s'", format)
		return
	}
}

func GetPublicIPV4(config models.IPSettings) (result string, err error) {
	if !config.Enabled {
		return
	}
	result, err = query(config.Query)
	ip := net.ParseIP(result)
	if ip.To4() == nil {
		err = fmt.Errorf("query ip is not ipv4 (%s)", result)
		result = ""
		return
	}

	return
}

func GetPublicIPV6(config models.IPSettings) (result string, err error) {
	if !config.Enabled {
		return
	}
	result, err = query(config.Query)
	ip := net.ParseIP(result)
	if ip.To4() != nil {
		err = fmt.Errorf("query ip is not ipv6 (%s)", result)
		result = ""
		return
	}

	return
}
