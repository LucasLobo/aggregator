package domain

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessEvents_WindowSize(t *testing.T) {
	input := []DurationEvent{
		{
			Timestamp: parseTime(t, "2018-12-26 18:11:08.509654"),
			Duration:  20,
		},
		{
			Timestamp: parseTime(t, "2018-12-26 18:15:19.903159"),
			Duration:  31,
		},
		{
			Timestamp: parseTime(t, "2018-12-26 18:23:19.903159"),
			Duration:  54,
		},
	}

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
			sw := NewSlidingWindow(tc.windowSize)

			expectedResults := createResultsByWindowSize(t, tc.windowSize)
			var actualResults []AverageDeliveryTime
			for _, event := range input {
				partial := sw.Process(event)
				actualResults = append(actualResults, partial...)

			}
			assert.Equal(t, expectedResults, actualResults)
		})
	}
}

func TestProcessEvents_OneEvent(t *testing.T) {
	sw := NewSlidingWindow(10)

	input := []DurationEvent{
		{
			Timestamp: parseTime(t, "2018-12-26 18:11:08.509654"),
			Duration:  20,
		},
	}
	expectedResults := []AverageDeliveryTime{
		{
			Date:                parseTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
	}

	var actualResults []AverageDeliveryTime
	for _, event := range input {
		partial := sw.Process(event)
		actualResults = append(actualResults, partial...)
	}

	assert.Equal(t, expectedResults, actualResults)
}

func createResultsByWindowSize(t *testing.T, windowSize int) []AverageDeliveryTime {
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

func createResultsWindowSize1(t *testing.T) []AverageDeliveryTime {
	return []AverageDeliveryTime{
		{
			Date:                parseTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 54,
		},
	}
}

func createResultsWindowSize10(t *testing.T) []AverageDeliveryTime {
	return []AverageDeliveryTime{
		{
			Date:                parseTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 31,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 42.5,
		},
	}
}

func createResultsWindowSize20(t *testing.T) []AverageDeliveryTime {
	return []AverageDeliveryTime{
		{
			Date:                parseTime(t, "2018-12-26 18:11:00.0000"),
			AverageDeliveryTime: 0,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:12:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:13:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:14:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:15:00.0000"),
			AverageDeliveryTime: 20,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:16:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:17:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:18:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:19:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:20:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:21:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:22:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:23:00.0000"),
			AverageDeliveryTime: 25.5,
		},
		{
			Date:                parseTime(t, "2018-12-26 18:24:00.0000"),
			AverageDeliveryTime: 35,
		},
	}
}

func parseTime(t *testing.T, val string) Time {
	var tt Time
	err := json.Unmarshal([]byte("\""+val+"\""), &tt)
	if err != nil {
		assert.FailNow(t, "Couldn't parse time", err)
	}
	return tt
}
