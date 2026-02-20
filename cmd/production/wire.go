//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/pdcgo/order_service/order_mutation"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/tracking_service"
	"github.com/pdcgo/tracking_service/thirdparties"
	"github.com/urfave/cli/v3"
)

func InitializeApp() (App, error) {
	wire.Build(
		http.NewServeMux,
		configs.NewProductionConfig,
		custom_connect.NewDefaultClientInterceptor,
		custom_connect.NewRegisterReflect,
		custom_connect.NewDefaultInterceptor,
		NewDatabase,
		NewCache,
		NewAuthorization,
		order_mutation.NewTagMutation,
		NewTrackingClient,
		// NewOrderClient,

		thirdparties.NewMultipleTracker,
		NewCheckOrder,
		tracking_service.NewRegister,
		NewApi,
		NewApp,
	)
	return &cli.Command{}, nil
}
