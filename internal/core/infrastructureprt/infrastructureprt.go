package infrastructureprt

import (
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type MovingAverageStorer interface {
	StoreMovingAverage([]domain.AverageDeliveryTime) error
}
