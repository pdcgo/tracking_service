package main

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/schema/services/tracking_iface/v1/tracking_ifaceconnect"
	"github.com/pdcgo/shared/pkg/cloud_logging"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"
)

type CheckOrderFunc cli.ActionFunc

func NewCheckOrder(
	db *gorm.DB,
	trackingClient tracking_ifaceconnect.TrackingServiceClient,
) CheckOrderFunc {

	return func(ctx context.Context, c *cli.Command) error {
		var err error

		if !c.Bool("log-local") {
			cloud_logging.SetCloudLoggingDefault()
		}

		// checking shipped
		err = ProcessShipped(ctx, db, trackingClient)
		if err != nil {
			return err
		}

		// checking return

		return nil

	}
}

type TrackItem struct {
	Receipt    string `json:"receipt"`
	ShippingID uint   `json:"shipping_id"`
	OrderID    uint   `json:"order_id"`
	TeamID     uint   `json:"team_id"`
	TxID       uint   `json:"tx_id"`
}

func ProcessShipped(
	ctx context.Context,
	db *gorm.DB,
	trackingClient tracking_ifaceconnect.TrackingServiceClient,
) error {
	// getting tracer
	// trace := otel.Tracer("")
	// spanCtx, span := trace.Start(ctx, "bulk_check_shipped")
	// defer span.End()

	// getting query
	query := db.
		Table("inv_transactions").
		Select([]string{
			"inv_transactions.id as tx_id",
			"inv_transactions.receipt as receipt",
			"orders.id as order_id",
			"orders.id as team_id",

			// "orders.order_ref_id as order_ref_id",
			// "orders.team_id as team_id",
			// "inv_transactions.warehouse_id",
			"inv_transactions.shipping_id",
		}).
		Joins("join orders on orders.invertory_tx_id = inv_transactions.id").
		// Where("shipping_id in ?", []uint{7, 8}).
		Where("orders.status in ?", []string{"courrier_shipped"})

	rows, err := query.Rows()
	if err != nil {
		// span.SetStatus(codes.Error, err.Error())
		return err
	}

	defer rows.Close()

	// iterate row
	for rows.Next() {

		var item TrackItem
		err = db.ScanRows(rows, &item)
		if err != nil {
			return err
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)

		_, err = trackingClient.TrackingProcess(ctxTimeout, &connect.Request[tracking_iface.TrackingProcessRequest]{
			Msg: &tracking_iface.TrackingProcessRequest{
				Payload: &tracking_iface.TrackingPayload{
					ShippingId: uint64(item.ShippingID),
					Receipt:    item.Receipt,
					Type: &tracking_iface.TrackingPayload_OrderShipped{
						OrderShipped: &tracking_iface.OrderShipped{
							OrderId: uint64(item.OrderID),
							TxId:    uint64(item.TxID),
							TeamId:  uint64(item.TeamID),
						},
					},
				},
			},
		})

		cancel()

		if err != nil {
			slog.Error(
				err.Error(),
				"order_id", item.OrderID,
				"receipt", item.Receipt,
				"tx_id", item.TxID,
				"team_id", item.TeamID,
			)
			continue
		}

		slog.Info("checked receipt", "receipt", item.Receipt)

	}

	// span.SetStatus(codes.Ok, "")
	return nil
}

func ProcessReturn() error {
	panic("unimplemented")
}

func ProcessRestock() error {
	panic("unimplemented")
}
