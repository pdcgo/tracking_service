package thirdparties

import (
	"fmt"

	"github.com/pdcgo/shared/configs"
	"github.com/pdcgo/shared/pkg/common_helper"
	"github.com/pdcgo/shared/pkg/raja_ongkir"
)

type MultipleTrackerConfig struct {
	ToniUndergroundEnpoint string
	RajaOngkirKey          []string
}

func NewMultipleTracker(cfg *configs.AppConfig) common_helper.NextFuncParam[*TrackProcess] {
	trackingCfg := cfg.TrackingService

	spxTracker := NewSpxUndergroundTracker()
	toniTracker := NewToniUndergroundTracker(trackingCfg.ToniUndergroundEnpoint)
	rajaOngkirTracker := NewRajaOngkirTracker(raja_ongkir.NewApiKey(trackingCfg.RajaOngkirKey))

	spxPriority := common_helper.NewChainParam(
		spxTracker,
		toniTracker,
		rajaOngkirTracker,
		fallbackErr,
	)

	defaultPriority := common_helper.NewChainParam(
		toniTracker,
		rajaOngkirTracker,
		fallbackErr,
	)

	return func(data *TrackProcess) (*TrackProcess, error) {
		var err error
		req := data.Req

		switch req.ShippingId {
		case 7, 8:
			_, err = spxPriority(data)
		case 39, 0:
			return data, fmt.Errorf("shipping id %d unsupported", req.ShippingId)
		default:
			_, err = defaultPriority(data)

		}

		if err != nil {
			return data, err
		}

		return data, nil
	}
}

func fallbackErr(next common_helper.NextFuncParam[*TrackProcess]) common_helper.NextFuncParam[*TrackProcess] {
	return func(data *TrackProcess) (*TrackProcess, error) {
		if data.Error != nil {
			return data, data.Error
		}
		return data, nil
	}
}
