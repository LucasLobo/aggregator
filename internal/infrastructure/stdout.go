package infrastructure

import (
	"encoding/json"
	"fmt"

	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

// StdOut is a simple implementation of a MovingAverageStorer that simply writes to the std output.
type StdOut struct {
}

func NewStdOut() StdOut {
	return StdOut{}
}

func (s StdOut) StoreMovingAverage(deliveryTimes []domain.AverageDeliveryTime) error {
	fmt.Println("OUTPUT START")
	for _, item := range deliveryTimes {
		bytes, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("error marshalling JSON: %w", err)
		}
		fmt.Println(string(bytes))
	}
	fmt.Println("OUTPUT END")
	return nil
}
