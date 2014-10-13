package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hamfist/artifacts-service/server"
)

var (
	// VersionString contains the git description
	VersionString = "?"

	// RevisionString contains the git revision
	RevisionString = "?"

	// GeneratedString contains the build date
	GeneratedString = "?"
)

func main() {
	log := logrus.New()
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s v=%s d=%s\n", c.App.Name, c.App.Version, GeneratedString)
	}

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
