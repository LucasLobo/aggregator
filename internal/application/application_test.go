package application

import (
	"github.com/lucaslobo/aggregator/internal/domain"
)

type mockStorer struct {
	store []domain.AverageDeliveryTime
}

func (ms *mockStorer) StoreMovingAverage(deliveryTime domain.AverageDeliveryTime) error {
	if ms.store == nil {
		ms.store = make([]domain.AverageDeliveryTime, 0)
	}
	ms.store = append(ms.store, deliveryTime)
	return nil
}

func (ms *mockStorer) StoreMovingAverageSlice(deliveryTimes []domain.AverageDeliveryTime) error {
	for _, deliveryItem := range deliveryTimes {
		err := ms.StoreMovingAverage(deliveryItem)
		if err != nil {
			return err
		}
	}
	return nil
}
