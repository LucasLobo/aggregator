package inboundprt

import (
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type MovingAverageCalculator interface {
	CalculateMovingAverage(events []domain.TranslationDelivered, windowSize int) error
}
