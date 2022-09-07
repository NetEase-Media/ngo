package ngo

import (
	"os"

	"github.com/NetEase-Media/ngo/pkg/config"
	"github.com/NetEase-Media/ngo/pkg/log"
	"github.com/urfave/cli/v2"
)

func Flag() {
	app := &cli.App{
		Name:    "ngo-app",
		Version: "v1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "app.yaml",
				Usage:   "the global config",
			},
			&cli.BoolFlag{
				Name:    "watch",
				Aliases: []string{"w"},
				Value:   false,
				Usage:   "watch config change event",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.New(c.String("config"), c.Bool("watch"))
			if err != nil {
				return err
			}
			config.SetConfiguration(cfg)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
