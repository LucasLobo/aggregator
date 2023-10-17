package inboundprt

import (
	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type MovingAverageCalculator interface {
	ProcessEvent(event domain.TranslationDelivered) error
}
