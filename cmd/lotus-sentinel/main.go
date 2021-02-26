package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/filecoin-project/lotus/lib/lotuslog"
)

var log = logging.Logger("sentinel")

func main() {
	lotuslog.SetupLogLevels()
	app := &cli.App{
		Name:                 "lotus-sentinel",
		Usage:                "filecoin blockchain monitoring and analysis",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "repo",
				EnvVars: []string{"SENTINEL_LOTUS_PATH"},
				Hidden:  true,
				Value:   "~/.lotus", // TODO: Consider XDG_DATA_HOME
			},
		},
		Commands: []*cli.Command{
			DaemonCmd,
			sentinelStartWatchCmd,
		},
	}
	app.Setup()
	lcli.RunApp(app)
}
