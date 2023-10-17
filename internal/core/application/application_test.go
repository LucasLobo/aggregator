package application

import (
	"fmt"
	"testing"
	"time"

	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockStorer struct {
	t     *testing.T
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

func (ms *mockStorer) Close() error {
	return nil
}

func TestProcessEvents_WindowSize(t *testing.T) {

	tests := []struct {
		windowSize int
	}{
		{
			windowSize: 1,
		},
		{
			windowSize: 10,
		},
		{
			windowSize: 20,
		},
		{
			windowSize: 9999,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("window size: %d", tc.windowSize), func(t *testing.T) {
			ms := mockStorer{
				t: t,
			}
			a := New(tc.windowSize, &ms)

			events := createEvents(t, 3)
			results := createResultsWindow(t, tc.windowSize)

			for _, event := range events {
				err := a.ProcessEvent(event)
				require.NoError(t, err)
			}

			assert.Equal(t, results, ms.store)
		})
	}
}

func TestProcessEvents_OneEvent(t *testing.T) {
	ms := mockStorer{
		t: t,
	}
	a := New(10, &ms)

	events := createEvents(t, 1)
	results := createResultsWindow(t, 10)[0:2]

	for _, event := range events {
		err := a.ProcessEvent(event)
		require.NoError(t, err)
	}

	assert.Equal(t, results, ms.store)
}

func mustGetTime(t *testing.T, val string) time.Time {

	tt, err := time.Parse("2006-01-02 15:04:05.999999", val)
	if err != nil {
		assert.FailNow(t, "Couldn't parse time", err)
	}

	return tt
}

func createEvents(t *testing.T, numberOfEvents int) []domain.TranslationDelivered {
	return []domain.TranslationDelivered{
		{
			Timestamp: mustGetTime(t, "2018-12-26 18:11:08.509654"),
			Duration:  20,
		},
		{
			Timestamp: mustGetTime(t, "2018-12-26 18:15:19.903159"),
			Duration:  31,
		},
		{
			Timestamp: mustGetTime(t, "2018-12-26 18:23:19.903159"),
			Duration:  54,
		},
	}[0:numberOfEvents]
}

func createResultsWindow(t *testing.T, windowSize int) []domain.AverageDeliveryTime {
	switch windowSize {
	case 1:
		return createResultsWindowSize1(t)
	case 10:
		return createResultsWindowSize10(t)
	case 20:
		return createResultsWindowSize20(t)
	case 9999:
		return createResultsWindowSize20(t)
	default:
		assert.FailNow(t, "unexpected window size")
		return nil
	}
}

func createResultsWindowSize1(t *testing.T) []domain.AverageDeliveryTime {
	return []domain.AverageDeliveryTime{
		{
			Date:                mustGetTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 54,
		},
	}
}

func createResultsWindowSize10(t *testing.T) []domain.AverageDeliveryTime {
	return []domain.AverageDeliveryTime{
		{
			Date:                mustGetTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 42.5,
		},
	}
}

func createResultsWindowSize20(t *testing.T) []domain.AverageDeliveryTime {
	return []domain.AverageDeliveryTime{
		{
			Date:                mustGetTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                mustGetTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 35,
		},
	}
}
