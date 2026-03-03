package thirdparties_test

import (
	"testing"

	"github.com/pdcgo/schema/services/tracking_iface/v1"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/shared/pkg/debugtool"
	"github.com/pdcgo/tracking_service/thirdparties"
	"github.com/stretchr/testify/assert"
)

func TestSpxUnderGround(t *testing.T) {

	tracker := common_helper.NewChainParam(
		thirdparties.NewSpxUndergroundTracker(),
	)

	data, err := tracker(&thirdparties.TrackProcess{
		Req: &tracking_iface.TrackingPayload{
			ShippingId: 7,
			Receipt:    "SPXID069027265211",
		},
		Res: &tracking_iface.TrackInfo{},
	})

	debugtool.LogJson(data)

	assert.Nil(t, err)
}
