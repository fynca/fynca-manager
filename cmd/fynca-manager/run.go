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
	"context"

	"git.underland.io/ehazlett/fynca/pkg/tracing"
	"github.com/ehazlett/fynca-manager/server"
	cli "github.com/urfave/cli/v2"
)

func run(clix *cli.Context) error {
	cfg := &server.Config{
		GRPCAddress:          clix.String("grpc-address"),
		TLSServerCertificate: clix.String("grpc-tls-server-cert"),
		TLSServerKey:         clix.String("grpc-tls-server-key"),
		TLSClientCertificate: clix.String("grpc-tls-client-cert"),
		TLSClientKey:         clix.String("grpc-tls-client-key"),
		APICORSDomain:        clix.String("api-cors-domain"),
		ListenAddr:           clix.String("listen-addr"),
		PublicDir:            clix.String("public-dir"),
		TraceEndpoint:        clix.String("trace-endpoint"),
		Environment:          clix.String("environment"),
	}

	tp, err := tracing.NewProvider(cfg.TraceEndpoint, "fynca-manager", cfg.Environment)
	if err != nil {
		return err
	}
	defer tp.Shutdown(context.Background())

	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	return srv.Run()
}
