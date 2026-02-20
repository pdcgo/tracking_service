package thirdparties

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/shared/pkg/spx_tracker"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewSpxUndergroundTracker() common_helper.NextHandlerParam[*TrackProcess] {
	var spxClient = spx_tracker.TrackClient{
		C: http.DefaultClient,
	}

	return func(next common_helper.NextFuncParam[*TrackProcess]) common_helper.NextFuncParam[*TrackProcess] {
		return func(data *TrackProcess) error {
			var err error
			req := data.Req
			trackInfo := data.Res
			trackInfo.Thirdparties = append(trackInfo.Thirdparties, tracking_iface.Thirdparty_THIRDPARTY_SPX_UNDERGROUND)

			res, err := spxClient.Track(req.Receipt)
			if err != nil {
				return next(data.WithError(err))
			}

			// debugtool.LogJson(data)

			trackInfo.Status, trackInfo.Histories, err = parseStatusAndHistory(req.Receipt, res)
			if err != nil {
				return next(data.WithError(err))
			}

			return nil
		}
	}
}

func parseStatusAndHistory(receipt string, data *spx_tracker.TrackResponse) (tracking_iface.Status, []*tracking_iface.HistoryItem, error) {
	var err error
	var history []*tracking_iface.HistoryItem
	var status tracking_iface.Status = tracking_iface.Status_STATUS_UNSPECIFIED

	// history = s.getHistory(data)
	switch data.Data.CurrentStatus {
	case spx_tracker.Delivered:
		status = tracking_iface.Status_STATUS_DELIVERED
	case spx_tracker.Lost:
		status = tracking_iface.Status_STATUS_LOST
	case spx_tracker.OnHold:
		status = tracking_iface.Status_STATUS_DELIVERY_FAILED
	case "Delivering":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "DOP_Received":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "SOC_Received":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "SOC_LHTransporting":
		// info.Status = db_models.DeliveryFailed
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "LMHub_Received":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "FMHub_Received":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "FMHub_LHTransporting":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS
	case "FMHub_Pickup_Done":
		status = tracking_iface.Status_STATUS_SHIPMENT_PROCESS

	case spx_tracker.Cancelled:
		status = tracking_iface.Status_STATUS_CANCEL

	case "Return_SOC_LHTransporting":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS
	case "Return_LMHub_Returned":
		status = tracking_iface.Status_STATUS_RETURNED
	case "Return_LMHub_LHTransporting":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS
	case "Return_LMHub_Received":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS

	case "Return_FMHub_Onhold":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS
	case "Return_SOC_Received":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS
	case "Return_FMHub_Received":
		status = tracking_iface.Status_STATUS_RETURN_PROCESS
	case "Return_FMHub_Returned":
		status = tracking_iface.Status_STATUS_RETURNED

	case "Created":
		status = tracking_iface.Status_STATUS_CREATED

	default:
		status = tracking_iface.Status_STATUS_UNSPECIFIED
		return status, history, fmt.Errorf("[spx] status %s not listed %s", data.Data.CurrentStatus, receipt)
	}

	for _, item := range data.Data.TrackingList {
		history = append(history, &tracking_iface.HistoryItem{
			Name: item.Status,
			Desc: item.Message,
			At:   timestamppb.New(time.Unix(item.Timestamp, 0)),
		})
	}

	return status, history, err

}
