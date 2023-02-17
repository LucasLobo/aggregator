package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
)

type assertStorer struct {
	t        *testing.T
	expected []domain.AverageDeliveryTime
}

func (as assertStorer) StoreMovingAverage(deliveryTimes []domain.AverageDeliveryTime) error {
	assert.Equal(as.t, as.expected, deliveryTimes)
	return nil
}

func Test_CalculateMovingAverage_HappyPath(t *testing.T) {
	results := createResults(t)
	domainResults := toAverageDeliveryTimeSlice(results)
	app := New(assertStorer{t: t, expected: domainResults})

	windowSize := 10
	events := createEvents(t)

	err := app.CalculateMovingAverage(events, windowSize)

	assert.NoError(t, err)
}

func Test_CalculateMovingAverage_NoEvents(t *testing.T) {
	app := New(assertStorer{t: t, expected: nil})

	windowSize := 10
	var events []domain.TranslationDelivered

	err := app.CalculateMovingAverage(events, windowSize)
	assert.EqualError(t, err, "no events provided")
}

func Test_calculateMovingAverage(t *testing.T) {

	var tests = []struct {
		name           string
		events         []domain.TranslationDelivered
		windowSize     int
		expectedResult []movingAverage
	}{
		{
			name:           "test case 1",
			events:         createEvents(t),
			windowSize:     10,
			expectedResult: createResults(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualRes := calculateMovingAverage(tt.events, tt.windowSize)
			for i := range actualRes {
				assert.Equal(t, tt.expectedResult[i].Timestamp, actualRes[i].Timestamp)
				assert.Equal(t, tt.expectedResult[i].AverageDuration, actualRes[i].AverageDuration)
			}
		})
	}
}

func mustGetTime(t *testing.T, val string) time.Time {

	tt, err := time.Parse("2006-01-02 15:04:05.999999", val)
	if err != nil {
		assert.FailNow(t, "Couldn't parse time", err)
	}

	return tt
}

func createEvents(t *testing.T) []domain.TranslationDelivered {
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
	}
}

func createResults(t *testing.T) []movingAverage {
	return []movingAverage{
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:11:00.0000"),
			AverageDuration: 0,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:12:00.0000"),
			AverageDuration: 20,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:13:00.0000"),
			AverageDuration: 20,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:14:00.0000"),
			AverageDuration: 20,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:15:00.0000"),
			AverageDuration: 20,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:16:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:17:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:18:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:19:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:20:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:21:00.0000"),
			AverageDuration: 25.5,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:22:00.0000"),
			AverageDuration: 31,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:23:00.0000"),
			AverageDuration: 31,
		},
		{
			Timestamp:       mustGetTime(t, "2018-12-26 18:24:00.0000"),
			AverageDuration: 42.5,
		},
	}
}
