package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/hamfist/artifacts-service/metadata"
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
	app.Usage = "CRUD FUN"
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
			Name:  "migrate",
			Usage: "run database migrations",
			Action: func(_ *cli.Context) {
				log.Debug("getting new server options")
				opts := server.NewOptions()

				log.WithFields(logrus.Fields{
					"opts": opts,
				}).Debug("spinning up database")

				db, err := metadata.NewDatabase(opts.DatabaseURL, log)
				if err != nil {
					log.Fatal(err)
				}

				log.Debug("migrating")

				err = db.Migrate(log)
				if err != nil {
					log.Fatal(err)
				}

				log.Info("database migration complete")
			},
		},
	}

	app.Run(os.Args)
}
