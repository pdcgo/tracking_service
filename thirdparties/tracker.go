package thirdparties

import (
	"github.com/pdcgo/schema/services/tracking_iface/v1"
)

type TrackProcess struct {
	Req   *tracking_iface.TrackingPayload
	Res   *tracking_iface.TrackInfo
	Error error
}

func (t *TrackProcess) WithError(err error) *TrackProcess {
	if t.Error != nil {
		return t
	}
	t.Error = err
	return t
}
