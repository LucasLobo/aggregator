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

func (s StdOut) StoreMovingAverage(item domain.AverageDeliveryTime) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	fmt.Println(string(bytes))
	return nil
}

func (s StdOut) Close() error {
	// there's no point in closing anything here, let's just return silently
	return nil
}
