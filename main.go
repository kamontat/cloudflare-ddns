package main

import (
	"os"

	"github.com/kamontat/cloudflare-ddns/cfw"
	"github.com/kamontat/cloudflare-ddns/models"
	"github.com/kc-workspace/go-lib/commandline"
	"github.com/kc-workspace/go-lib/commandline/commands"
	m "github.com/kc-workspace/go-lib/commandline/models"
	"github.com/kc-workspace/go-lib/commandline/plugins"
	"github.com/kc-workspace/go-lib/logger"
)

// assign from goreleaser
var (
	short   string = "cfddns"
	name    string = "cloudflare-ddns"
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
	builtBy string = "manually"
)

func main() {
	var err = commandline.New(&m.Metadata{
		Short:   short,
		Name:    name,
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
		Usage: `
Update cloudflare record
`,
	}).
		Plugin(plugins.SupportHelp()).
		Plugin(plugins.SupportVersion()).
		Plugin(plugins.SupportLogLevel(logger.INFO)).
		Plugin(plugins.SupportConfig([]string{"{{.current}}/configs"})).
		Plugin(plugins.SupportDotEnv(false)).
		Plugin(plugins.SupportVar()).
		Command(&commands.Command{
			Name:    commands.DEFAULT,
			Aliases: []string{"update"},
			Executor: func(p *commands.ExecutorParameter) error {
				var config, err = models.NewConfig(p.Config)
				if err != nil {
					return err
				}

				wrapper, err := cfw.New(config, p.Logger.New())
				if err != nil {
					return err
				}

				wrapper.UpsertRecords()
				return nil
			},
		}).Start(os.Args)

	if err != nil {
		panic(err)
	}
}
