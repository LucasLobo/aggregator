package inboundprt

import (
	"github.com/lucaslobo/aggregator/internal/core/domain"
)

type MovingAverageCalculator interface {
	ProcessEvent(event domain.TranslationDelivered) error
}
