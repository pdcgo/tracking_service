package tracking_service

import (
	"net/http"

	"github.com/pdcgo/schema/services/tracking_iface/v1/tracking_ifaceconnect"
	"github.com/pdcgo/shared/custom_connect"
	"github.com/pdcgo/shared/interfaces/authorization_iface"
	"github.com/pdcgo/shared/interfaces/order_iface"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/tracking_service/thirdparties"
	"github.com/pdcgo/tracking_service/tracking"
	"gorm.io/gorm"
)

type ServiceReflectNames []string
type RegisterHandler func() ServiceReflectNames

func NewRegister(
	mux *http.ServeMux,
	db *gorm.DB,
	auth authorization_iface.Authorization,
	tracker common_helper.NextFuncParam[*thirdparties.TrackProcess],
	defaultInterceptor custom_connect.DefaultInterceptor,
	tagMutation order_iface.OrderTagMutation,
) RegisterHandler {

	return func() ServiceReflectNames {
		grpcReflects := ServiceReflectNames{}

		path, handler := tracking_ifaceconnect.NewTrackingServiceHandler(
			tracking.NewTrackingService(db, tracker, tagMutation),
			defaultInterceptor,
		)
		mux.Handle(path, handler)
		grpcReflects = append(grpcReflects, tracking_ifaceconnect.TrackingServiceName)

		return grpcReflects
	}
}
