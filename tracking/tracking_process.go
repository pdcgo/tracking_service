package tracking

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/tracking_service/thirdparties"
)

// TrackingProcess implements [tracking_ifaceconnect.TrackingServiceHandler].
func (t *trackingServiceImpl) TrackingProcess(
	ctx context.Context,
	req *connect.Request[tracking_iface.TrackingProcessRequest],

) (*connect.Response[tracking_iface.TrackingProcessResponse], error) {
	var err error
	pay := req.Msg
	res := tracking_iface.TrackingProcessResponse{
		TrackInfo: &tracking_iface.TrackInfo{
			Histories:    []*tracking_iface.HistoryItem{},
			Thirdparties: []tracking_iface.Thirdparty{},
		},
	}

	trackProcess := &thirdparties.TrackProcess{Req: pay.Payload, Res: res.TrackInfo, Error: nil}

	err = common_helper.NewChain(
		func(next common_helper.NextFunc) common_helper.NextFunc {
			return func() error { // checking tracking
				err := t.tracker(trackProcess)
				if err != nil {
					return err
				}

				return next()
			}
		},
		func(next common_helper.NextFunc) common_helper.NextFunc {
			return func() error { // removing tag order
				switch extra := pay.Payload.Type.(type) {
				case *tracking_iface.TrackingPayload_OrderShipped:

					err = t.tagMutation.RemoveAllFrom(db_models.RelationFromTracking, []uint{
						uint(extra.OrderShipped.OrderId),
					})

					if err != nil {
						return err
					}
				default:
					return errors.New("type data in payload not unsupported")
				}

				return next()
			}
		},

		func(next common_helper.NextFunc) common_helper.NextFunc {
			return func() error { // tagging order and getting order dan transaksi

				var err error

				var ord db_models.Order
				var tx db_models.InvTransaction
				var orderID uint64
				info := res.TrackInfo

				switch extra := pay.Payload.Type.(type) {
				case *tracking_iface.TrackingPayload_OrderShipped:
					err = t.
						db.
						Model(&db_models.Order{}).
						Select(
							"id",
							"status",
						).
						First(&ord, extra.OrderShipped.OrderId).
						Error

					if err != nil {
						return err
					}

					err = t.
						db.
						Model(&db_models.InvTransaction{}).
						Select(
							"id",
							"status",
							"receipt",
							"type",
						).
						First(&tx, extra.OrderShipped.TxId).
						Error

					if err != nil {
						return err
					}

					orderID = extra.OrderShipped.OrderId

				default:
					return errors.New("type data in payload not unsupported")
				}

				tags := []string{}

				switch tx.Type {
				case db_models.InvTxOrder:
					switch ord.Status {
					case db_models.OrdCompleted:
						switch info.Status {
						case tracking_iface.Status_STATUS_LOST:
							tags = append(tags, "hilang")
						case tracking_iface.Status_STATUS_DELIVERY_FAILED:
							tags = append(tags, "pengiriman_gagal")
						case tracking_iface.Status_STATUS_CANCEL:
							tags = append(tags, "batal")
						case tracking_iface.Status_STATUS_RETURN_PROCESS:
							tags = append(tags, "return_ke_gudang")
						case tracking_iface.Status_STATUS_RETURNED:
							tags = append(tags, "return_sampai_gudang")
						}
					case db_models.OrdShipped,
						db_models.OrdReadyForPacking,
						db_models.OrdReadyForCourrier,
						db_models.OrdCreated,
						db_models.OrdProcess,
						db_models.OrdProductPick,
						db_models.OrdProblem,
						db_models.OrdCourrierShipped:
						switch info.Status {
						case tracking_iface.Status_STATUS_DELIVERED:
							tags = append(tags, "sampai")
						case tracking_iface.Status_STATUS_LOST:
							tags = append(tags, "hilang")
						case tracking_iface.Status_STATUS_DELIVERY_FAILED:
							tags = append(tags, "pengiriman_gagal")
						case tracking_iface.Status_STATUS_CANCEL:
							tags = append(tags, "batal")
						case tracking_iface.Status_STATUS_RETURN_PROCESS:
							tags = append(tags, "return_ke_gudang")
						case tracking_iface.Status_STATUS_RETURNED:
							tags = append(tags, "return_sampai")
						}
					case db_models.OrdCancel,
						db_models.OrdReturnCompleted:
					}

				case db_models.InvTxReturn:
					switch ord.Status {
					case db_models.OrdReturn, db_models.OrdReturnProblem:
						switch info.Status {
						case tracking_iface.Status_STATUS_DELIVERED:
							tags = append(tags, "return_sampai")
						case tracking_iface.Status_STATUS_LOST:
							tags = append(tags, "return_hilang")
						case tracking_iface.Status_STATUS_DELIVERY_FAILED:
							tags = append(tags, "return_gagal")
						case tracking_iface.Status_STATUS_CANCEL:
							tags = append(tags, "return_batal")
						case tracking_iface.Status_STATUS_RETURNED:
							tags = append(tags, "return_sampai")
						}
					}
				}

				err = t.tagMutation.Add(
					db_models.RelationFromTracking,
					[]uint{uint(orderID)},
					tags,
				)

				if err != nil {
					return err
				}

				return next()
			}
		},
	)

	if err != nil {
		return connect.NewResponse(&res), err
	}

	return connect.NewResponse(&res), nil
}
