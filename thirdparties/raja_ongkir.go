package thirdparties

import (
	"fmt"
	"time"

	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/shared/pkg/raja_ongkir"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RajaOngkirKey struct {
	Keys []string
}

func NewRajaOngkirTracker(keys *raja_ongkir.ApiKey) common_helper.NextHandlerParam[*TrackProcess] {

	return func(next common_helper.NextFuncParam[*TrackProcess]) common_helper.NextFuncParam[*TrackProcess] {
		return func(data *TrackProcess) (*TrackProcess, error) {
			req := data.Req
			trackInfo := data.Res

			// setting thirdparty
			trackInfo.Thirdparties = append(trackInfo.Thirdparties, tracking_iface.Thirdparty_THIRDPARTY_RAJAONGKIR)

			// getting courier
			courier := IDHelper(req.ShippingId)
			if courier == "" {
				return data, fmt.Errorf("unsupported tracking with shipping id %d", req.ShippingId)
			}

			res, err := raja_ongkir.KomerceTrack(keys, req.Receipt, courier)
			if err != nil {
				return next(
					data.WithError(err),
				)
			}

			// getting history
			histories := []*tracking_iface.HistoryItem{}
			for _, item := range res.Data.Manifest {
				ts, err := item.GetTimestamp()
				if err != nil {
					return data, err
				}

				histories = append(histories, &tracking_iface.HistoryItem{
					Name: item.CityName,
					Desc: item.ManifestDescription,
					At:   timestamppb.New(time.Unix(ts, 0)),
				})

			}
			trackInfo.Histories = histories

			switch res.Data.Summary.Status {
			case raja_ongkir.Delivered:
				trackInfo.Status = tracking_iface.Status_STATUS_DELIVERED
			case raja_ongkir.OnProcess:
				trackInfo.Status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
			case raja_ongkir.ReturnProcess:
				trackInfo.Status = tracking_iface.Status_STATUS_RETURN_PROCESS

			default:
				trackInfo.Status = tracking_iface.Status_STATUS_UNSPECIFIED
				return data, fmt.Errorf("[raja_ongkir] status %s not listed %s", res.Data.Summary.Status, req.Receipt)

			}

			return data, nil
		}
	}
}

func IDHelper(id uint64) string {
	switch id {
	case 9, 10, 11:
		return "anteraja"
	case 12, 13, 14, 15:
		return "sicepat"
	case 16:
		return "ninja"
	case 17, 18, 19, 20:
		return "jnt"
	case 21, 22, 23, 24:
		return "jne"
	case 25:
		return "ide"
	case 26:
		return "pos"
	case 28:
		return "rex"
	case 35:
		return "wahana"
	default:
		return ""
	}
}
