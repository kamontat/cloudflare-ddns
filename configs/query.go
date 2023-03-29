package configs

import (
	"fmt"
	"io"
	"strings"

	"github.com/kamontat/cloudflare-ddns/clients"
	"github.com/kc-workspace/go-lib/configs"
	"github.com/kc-workspace/go-lib/mapper"
)

type Query struct {
	Url       *string `json:"url" yaml:"url"`
	Format    *string `json:"format" yaml:"format"`
	Separator *string `json:"separator" yaml:"separator"`
	Key       *string `json:"key" yaml:"key"`
}

func (q *Query) GetURL() string {
	if q.Url == nil {
		return ""
	}
	return *q.Url
}

func (q *Query) GetFormat() string {
	if q.Format == nil {
		return ""
	}
	return *q.Format
}

func (q *Query) GetSeparator() string {
	if q.Separator == nil {
		return ""
	}
	return *q.Separator
}

func (q *Query) GetKey() string {
	if q.Key == nil {
		return ""
	}
	return *q.Key
}

func (q *Query) Query() (result string, err error) {
	if q.GetURL() == "" {
		err = fmt.Errorf("query.url is required")
		return
	}
	if q.GetFormat() == "" {
		err = fmt.Errorf("query.format is required")
		return
	}

	req, err := clients.Default.Get(q.GetURL())
	if err != nil {
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	return q.convert(body)
}

func (q *Query) convert(body []byte) (result string, err error) {
	switch q.GetFormat() {
	case "text":
		result = string(body)
		return
	case "json":
		var m mapper.Mapper
		m, err = mapper.FromJson(body)
		if err != nil {
			return
		}

		result, err = m.Se(q.GetKey())
		return
	case "kv":
		var lines = strings.Split(string(body), "\n")
		for _, line := range lines {
			var k, v, ok = configs.ParseOverride(strings.Replace(line, q.GetSeparator(), "=", 1))
			if ok && q.GetKey() == k {
				result = strings.TrimSpace(v)
				return
			}
		}

		err = fmt.Errorf("cannot found key '%s' in response", q.GetKey())
		return
	default:
		err = fmt.Errorf("invalid format name '%s'", q.GetFormat())
		return
	}
}
