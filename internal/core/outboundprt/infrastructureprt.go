package outboundprt

import (
	"github.com/lucaslobo/aggregator/internal/core/domain"
)

type MovingAverageStorer interface {
	// StoreMovingAverage stores one domain.AverageDeliveryTime
	StoreMovingAverage(domain.AverageDeliveryTime) error

	// StoreMovingAverageSlice stores a slice of domain.AverageDeliveryTime
	StoreMovingAverageSlice([]domain.AverageDeliveryTime) error

	// Close closes the underlying resource/connection of the MovingAverageStorer
	Close() error
}
