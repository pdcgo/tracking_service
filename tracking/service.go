package tracking

import (
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/tracking_service/thirdparties"
)

type trackingServiceImpl struct {
	tracker common_helper.NextFuncParam[*thirdparties.TrackProcess]
}

func NewTrackingService(tracker common_helper.NextFuncParam[*thirdparties.TrackProcess]) *trackingServiceImpl {
	return &trackingServiceImpl{tracker}
}
