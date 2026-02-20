package tracking

import (
	"context"

	"connectrpc.com/connect"
	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/tracking_service/thirdparties"
)

// TrackingGet implements [tracking_ifaceconnect.TrackingServiceHandler].
func (t *trackingServiceImpl) TrackingGet(
	ctx context.Context,
	req *connect.Request[tracking_iface.TrackingGetRequest],
) (*connect.Response[tracking_iface.TrackingGetResponse], error) {

	var err error
	pay := req.Msg
	res := tracking_iface.TrackingGetResponse{
		TrackInfo: &tracking_iface.TrackInfo{
			Histories:    []*tracking_iface.HistoryItem{},
			Thirdparties: []tracking_iface.Thirdparty{},
		},
	}

	trackProcess := &thirdparties.TrackProcess{Req: pay.Payload, Res: res.TrackInfo, Error: nil}

	// err = common_helper.NewChain(
	// 	func(next common_helper.NextFunc) common_helper.NextFunc {
	// 		return func() error { // checking tracking
	// 			err := t.tracker(trackProcess)
	// 			if err != nil {
	// 				return err
	// 			}

	// 			return next()
	// 		}
	// 	},
	// 	func(next common_helper.NextFunc) common_helper.NextFunc {
	// 		return func() error { // removing tag order
	// 			return next()
	// 		}
	// 	},
	// 	func(next common_helper.NextFunc) common_helper.NextFunc {
	// 		return func() error { // filtering jika tidak ada order
	// 			return next()
	// 		}
	// 	},
	// 	func(next common_helper.NextFunc) common_helper.NextFunc {
	// 		return func() error { // getting order dan transaksi
	// 			return next()
	// 		}
	// 	},
	// 	func(next common_helper.NextFunc) common_helper.NextFunc {
	// 		return func() error { // tagging order
	// 			return next()
	// 		}
	// 	},
	// )

	err = t.tracker(trackProcess)

	if err != nil {
		return connect.NewResponse(&res), err
	}

	return connect.NewResponse(&res), nil

}
