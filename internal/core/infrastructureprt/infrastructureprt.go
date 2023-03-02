package infrastructureprt

import (
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type MovingAverageStorer interface {
	// StoreMovingAverage stores one domain.AverageDeliveryTime
	StoreMovingAverage(domain.AverageDeliveryTime) error
	StoreMovingAverageSlice([]domain.AverageDeliveryTime) error

	// Close allows us to close the underlying resource/connection of the MovingAverageStorer
	Close() error
}
