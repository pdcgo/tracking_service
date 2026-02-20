package tracking

import (
	"github.com/pdcgo/shared/interfaces/order_iface"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/tracking_service/thirdparties"
	"gorm.io/gorm"
)

type trackingServiceImpl struct {
	db          *gorm.DB
	tracker     common_helper.NextFuncParam[*thirdparties.TrackProcess]
	tagMutation order_iface.OrderTagMutation
}

func NewTrackingService(
	db *gorm.DB,
	tracker common_helper.NextFuncParam[*thirdparties.TrackProcess],
	tagMutation order_iface.OrderTagMutation,
) *trackingServiceImpl {
	return &trackingServiceImpl{db, tracker, tagMutation}
}
