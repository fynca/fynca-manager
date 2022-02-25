// Copyright 2022 Evan Hazlett
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	app.Usage = "Fynca Manager"
	app.Version = version.Version
	app.Authors = []*cli.Author{
		{
			Name: "@ehazlett",
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
		&cli.StringFlag{
			Name:  "trace-endpoint",
			Usage: "endpoint to send trace data",
		},
		&cli.StringFlag{
			Name:  "environment",
			Usage: "environment name for trace identification",
			Value: "dev",
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
