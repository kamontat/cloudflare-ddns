package cfw

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/configs"
	"github.com/kc-workspace/go-lib/logger"
	"github.com/kc-workspace/go-lib/mapper"
	"github.com/kc-workspace/go-lib/utils"
)

func GetPublicIP(config models.IPSettings) (*net.IP, error) {
	if !config.Enabled {
		return nil, nil
	}

	var result, err = query(config.Query)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(result)
	return &ip, nil
}

func GetFullDomain(name string, zone string) string {
	if name == "" || name == "@" || name == "." {
		return zone
	}
	return utils.JoinString(".", name, zone)
}

func GetTTL(input string, log *logger.Logger) (ttl int) {
	var min = 30
	var max = 86400

	ttl = 1 // automatic ttl
	if input != "" {
		// 60 and 86400
		duration, err := time.ParseDuration(input)
		if err != nil {
			log.Warnf("cannot parse ttl duration, fallback to 'automatic'")
			return
		}

		ttl = int(duration.Seconds())
		if ttl < min {
			log.Warnf("minimum ttl is %d, force ttl to be %d", min, min)
			ttl = min
		} else if ttl > max {
			log.Warnf("maximum ttl is %d, force ttl to be %d", max, max)
			ttl = max
		}
	}

	return
}

func query(query models.IPQuerySettings) (result string, err error) {
	if query.Url == "" {
		err = errors.New("query.url is required")
		return
	}
	if query.Format == "" {
		err = errors.New("query.format is required")
		return
	}

	req, err := DefaultClient.Get(query.Url)
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
				result = strings.TrimSpace(v)
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
