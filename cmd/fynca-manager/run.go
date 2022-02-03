package main

import (
	"github.com/ehazlett/fynca-manager/server"
	cli "github.com/urfave/cli/v2"
)

func run(cx *cli.Context) error {
	cfg := &server.Config{
		GRPCAddress:          cx.String("grpc-address"),
		TLSServerCertificate: cx.String("grpc-tls-server-cert"),
		TLSServerKey:         cx.String("grpc-tls-server-key"),
		TLSClientCertificate: cx.String("grpc-tls-client-cert"),
		TLSClientKey:         cx.String("grpc-tls-client-key"),
		APICORSDomain:        cx.String("api-cors-domain"),
		ListenAddr:           cx.String("listen-addr"),
		PublicDir:            cx.String("public-dir"),
	}
	srv, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	return srv.Run()
}
