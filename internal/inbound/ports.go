package inbound

import "github.com/lucaslobo/aggregator/internal/domain"

type MovingAverageCalculator interface {
	ProcessEvent(event domain.DurationEvent) error
}
