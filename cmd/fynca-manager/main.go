package main

import (
	"os"

	"github.com/ehazlett/fynca-manager/version"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = version.Name
	app.Usage = "Antimatter Cosmos Server"
	app.Version = version.Version
	app.Authors = []*cli.Author{
		{
			Name: "Project Antimatter Team",
		},
	}
	app.Before = func(c *cli.Context) error {
		// enable debug
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("debug enabled")
		}
		return nil
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"D"},
			Usage:   "enable debug",
		},
		&cli.StringFlag{
			Name:    "grpc-address",
			Aliases: []string{"a"},
			Usage:   "grpc address to fynca server",
			Value:   "127.0.0.1:8080",
		},
		&cli.StringFlag{
			Name:  "grpc-tls-server-cert",
			Usage: "grpc address tls server certificate",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "grpc-tls-server-key",
			Usage: "grpc address tls server key",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "grpc-tls-client-cert",
			Usage: "grpc address tls client certificate",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "grpc-tls-client-key",
			Usage: "grpc address tls client key",
			Value: "",
		},
		&cli.StringFlag{
			Name:    "listen-addr",
			Aliases: []string{"l"},
			Usage:   "listen address for server",
			Value:   ":9080",
		},
		&cli.StringFlag{
			Name:  "api-cors-domain",
			Usage: "domain for cors requests",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "public-dir",
			Usage: "directory that contains the Fynca Manager UI web application",
			Value: "./public",
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
