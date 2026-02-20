package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/tracking_service"
	"github.com/urfave/cli/v3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ApiFunc cli.ActionFunc

func NewApi(
	mux *http.ServeMux,
	reflectorRegister custom_connect.RegisterReflectFunc,
	trackingService tracking_service.RegisterHandler,

) ApiFunc {
	return func(ctx context.Context, c *cli.Command) error {

		// registering api
		grpcReflectNames := []string{}
		grpcReflectNames = append(grpcReflectNames, trackingService()...)

		// registering reflector
		reflectorRegister(grpcReflectNames)

		// running api
		port := os.Getenv("PORT")
		if port == "" {
			port = "8085"
		}

		host := os.Getenv("HOST")
		listen := fmt.Sprintf("%s:%s", host, port)
		log.Println("listening on", listen)

		return http.ListenAndServe(
			listen,
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(
				custom_connect.WithCORS(mux),
				&http2.Server{}),
		)

	}
}
