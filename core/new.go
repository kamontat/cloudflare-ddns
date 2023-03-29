package core

import (
	"github.com/kamontat/cloudflare-ddns/cloudflare"
	"github.com/kamontat/cloudflare-ddns/configs"
	"github.com/kc-workspace/go-lib/commandline/commands"
)

func New(p *commands.ExecutorParameter) (core *Core, err error) {
	config, err := configs.New(p.Config)
	if err != nil {
		return
	}

	cf, err := cloudflare.New(config.Secrets)
	if err != nil {
		return
	}

	core = &Core{
		meta:       p.Meta,
		cloudflare: cf,
		config:     config,
		logger:     p.Logger,
	}

	return
}
