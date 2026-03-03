package main

import (
	"context"
	"net/http"
	"os"

	"github.com/pdcgo/schema/services/order_iface/v1/order_ifaceconnect"
	"github.com/pdcgo/schema/services/tracking_iface/v1/tracking_ifaceconnect"
	"github.com/pdcgo/shared/authorization"
	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/db_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/pkg/ware_cache"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"
)

func NewOrderClient(
	cfg *configs.AppConfig,
	defaultClientInterceptor custom_connect.DefaultClientInterceptor,
) order_ifaceconnect.OrderServiceClient {
	return order_ifaceconnect.NewOrderServiceClient(
		http.DefaultClient,
		cfg.OrderService.Endpoint,
	)
}

func NewTrackingClient(
	cfg *configs.AppConfig,
	defaultClientInterceptor custom_connect.DefaultClientInterceptor,
) tracking_ifaceconnect.TrackingServiceClient {
	return tracking_ifaceconnect.NewTrackingServiceClient(
		http.DefaultClient,
		cfg.TrackingService.Endpoint,
	)
}

func NewDatabase(
	cfg *configs.AppConfig,
) (*gorm.DB, error) {
	return db_connect.NewProductionDatabase("tracking-service", &cfg.Database)
}

func NewCache(cfg *configs.AppConfig) (ware_cache.Cache, error) {
	return ware_cache.NewCustomCache(
		cfg.CacheService.Endpoint,
		// "http://localhost:8080",
	), nil
}

func NewAuthorization(
	cfg *configs.AppConfig,
	db *gorm.DB,
	cache ware_cache.Cache,
) authorization_iface.Authorization {
	return authorization.NewAuthorization(cache, db, cfg.JwtSecret)
}

type App *cli.Command

func NewApp(
	checkOrder CheckOrderFunc,
	api ApiFunc,
) App {
	return &cli.Command{
		Name: "tracking service",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "log-local",
				Value: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "check_order",
				Action: cli.ActionFunc(checkOrder),
			},
		},
		Action: cli.ActionFunc(api),
	}
}

func main() {

	cancel, err := custom_connect.InitTracer("tracking-service")
	if err != nil {
		panic(err)
	}
	defer cancel(context.Background())

	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	var command *cli.Command = app

	err = command.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
