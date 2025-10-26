package application

import (
	"github.com/lucaslobo/aggregator/internal/domain"
)

type Application struct {
	storer storer
	sw     *domain.SlidingWindow
}

type storer interface {
	StoreMovingAverage(domain.AverageDeliveryTime) error
	StoreMovingAverageSlice([]domain.AverageDeliveryTime) error
}

func New(windowSize int, storer storer) *Application {
	return &Application{
		storer: storer,
		sw:     domain.NewSlidingWindow(windowSize),
	}
}

// ProcessEvent calculates the moving average for all time-buckets since the last event.
// The moving-average is calculated based on the windowSize provided in the Init method
func (a *Application) ProcessEvent(event domain.DurationEvent) error {
	avgDeliveryTime := a.sw.Process(event)
	err := a.storer.StoreMovingAverage(avgDeliveryTime)
	if err != nil {
		return err
	}
	return nil
}
