package thirdparties

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/pkg/common_helper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ToniRes struct {
	Message string              `json:"message"`
	Data    db_models.TrackInfo `json:"data"`
}

func NewToniUndergroundTracker(endpoint string) common_helper.NextHandlerParam[*TrackProcess] {
	return func(next common_helper.NextFuncParam[*TrackProcess]) common_helper.NextFuncParam[*TrackProcess] {
		return func(data *TrackProcess) error {

			var hasil ToniRes
			var err error

			req := data.Req
			trackInfo := data.Res.TrackInfo

			// setting thirdparty
			trackInfo.Thirdparties = append(trackInfo.Thirdparties, tracking_iface.Thirdparty_THIRDPARTY_TONI_UNDERGROUND)

			params := url.Values{}
			params.Add("awb", req.Receipt)
			params.Add("s_id", strconv.Itoa(int(req.ShippingId)))
			uri := fmt.Sprintf("%s/track_awb?%s", endpoint, params.Encode())

			res, err := http.Get(uri)
			if err != nil {
				return next(data.WithError(err))
			}

			if res.StatusCode != 200 {
				dd, _ := io.ReadAll(res.Body)
				err = fmt.Errorf("[server toni] %s", dd)
				return next(data.WithError(err))
			}

			err = json.NewDecoder(res.Body).Decode(&hasil)
			if err != nil {
				return next(data.WithError(err))
			}

			if hasil.Message != "" {
				err = errors.New(hasil.Message)
				return next(data.WithError(err))
			}

			// setting history
			histories := []*tracking_iface.HistoryItem{}
			for _, item := range hasil.Data.History {
				ts := time.Unix(item.Timestamp, 0)
				histories = append(histories, &tracking_iface.HistoryItem{
					Name: item.Name,
					Desc: item.Desc,
					At:   timestamppb.New(ts),
				})
			}
			trackInfo.Histories = histories

			// setting status
			switch hasil.Data.Status {
			case db_models.Created:
				trackInfo.Status = tracking_iface.Status_STATUS_CREATED
			case db_models.DeliveryFailed:
				trackInfo.Status = tracking_iface.Status_STATUS_DELIVERY_FAILED
			case db_models.ShipmentProcess:
				trackInfo.Status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
			case db_models.TrackCancel:
				trackInfo.Status = tracking_iface.Status_STATUS_CANCEL
			case db_models.TrackLost:
				trackInfo.Status = tracking_iface.Status_STATUS_LOST
			case db_models.TrackReturnProcess:
				trackInfo.Status = tracking_iface.Status_STATUS_RETURN_PROCESS
			case db_models.TrackReturned:
				trackInfo.Status = tracking_iface.Status_STATUS_RETURNED
			default:
				trackInfo.Status = tracking_iface.Status_STATUS_UNSPECIFIED

			}

			return nil
		}
	}
}
