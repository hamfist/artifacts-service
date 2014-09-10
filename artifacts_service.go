package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hamfist/artifacts-service/server"
)

var (
	// VersionString contains the git description
	VersionString = "?"
)

func main() {
	log := logrus.New()
	app := cli.NewApp()
	app.Name = "artifacts-service"
	app.Version = VersionString
	app.Commands = []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "run the HTTP thing",
			Action: func(_ *cli.Context) {
				server.Main(log)
			},
		},
		{
			Name:      "migrate",
			ShortName: "m",
			Usage:     "run database migrations",
			Action: func(_ *cli.Context) {
				server.MigratorMain(log)
			},
		},
	}

	app.Run(os.Args)
}
