package application

import (
	"errors"
	"time"

	"github.com/lucaslobo/aggregator-cli/internal/core/domain"
	"github.com/lucaslobo/aggregator-cli/internal/core/infrastructureprt"
)

type Application struct {
	storer infrastructureprt.MovingAverageStorer
}

func New(storer infrastructureprt.MovingAverageStorer) Application {
	return Application{
		storer: storer,
	}
}

func (a Application) CalculateMovingAverage(events []domain.TranslationDelivered, windowSize int) error {
	items := calculateMovingAverage(events, windowSize)

	if items == nil {
		return errors.New("no events provided")
	}

	domainItems := toAverageDeliveryTimeSlice(items)
	return a.storer.StoreMovingAverage(domainItems)
}

// movingAverage represents the moving have in a single timestamp
type movingAverage struct {
	Timestamp       time.Time
	AverageDuration float32
}

type state struct {
	count    int
	duration int
}

// calculateMovingAverage calculates the moving average delivery time of the events with the given windowSize starting
// from the earliest event's timestamp until the latest event's timestamp.
// Each event's duration counts for the average for the minutes following timestamp up to the windowSize.
// e.g., windowSize = 2, 2018-12-26 18:11:08.509654 counts towards 2018-12-26 18:12:00 and 2018-12-26 18:13:00
func calculateMovingAverage(events []domain.TranslationDelivered, windowSize int) []movingAverage {
	if len(events) == 0 {
		return nil
	}
	// NOTE: I've tried to implement this sliding window average in multiple ways, including a slidingWindow struct
	// with the main method and multiple sub-methods. e.g., advanceHead, advanceTail, calculateAverage.
	// Usually for longer functions like this one I try to split it into multiple smaller functions that have a clear
	// responsibility, but I found that in this particular case it only made the code more confusing and complex.
	// In the end, I considered that the algorithm is the most readable when it's all contained in the same function.

	start := events[0].Timestamp.Truncate(time.Minute)
	end := events[len(events)-1].Timestamp.Truncate(time.Minute).Add(time.Minute).Add(time.Minute)
	duration := calculateDurationInMinutes(start, end) + 1

	// perMinuteAverage contains the moving average for each minute bucket
	perMinuteAverage := make([]movingAverage, 0, duration)

	// bucketState contains the partial state for each bucket in the running average
	bucketState := map[time.Time]state{}

	// runningState contains the current complete state of the sliding window
	runningState := state{}

	head := start
	tail := start
	currentEventIndex := 0

	// slide the window until we reach the end
	for head.Before(end) {

		currentEventDuration := 0
		currentEventCount := 0

		// we must count the contributions of the current event to the current timestamp
		// as long as we have not processed all events already
		if currentEventIndex < len(events) {

			currentEvent := events[currentEventIndex]
			eventTime := currentEvent.Timestamp.Truncate(time.Minute).Add(time.Minute)

			// if the event matches with the current timestamp, then we count it for our sliding window
			if eventTime.Equal(head) {
				currentEventDuration = currentEvent.Duration
				currentEventCount = 1
				currentEventIndex += 1

				runningState.count += 1
				runningState.duration += currentEventDuration
			}
		}

		bucketState[head] = state{
			duration: currentEventDuration,
			count:    currentEventCount,
		}

		// if we have more values than the window size, it means we must remove the first element from the tail
		// and move the tail forwards
		if len(bucketState) > windowSize {
			runningState.count -= bucketState[tail].count
			runningState.duration -= bucketState[tail].duration

			delete(bucketState, tail)
			tail = tail.Add(time.Minute)
		}

		// once we're done, we calculate the average for the current position
		average := float32(0)
		if runningState.count != 0 {
			average = float32(runningState.duration) / float32(runningState.count)
		}

		// and append the value to the perMinuteAverage
		perMinuteAverage = append(perMinuteAverage, movingAverage{
			Timestamp:       head,
			AverageDuration: average,
		})

		// finally we advance the head
		head = head.Add(time.Minute)
	}

	return perMinuteAverage
}

func calculateDurationInMinutes(startTime, endTime time.Time) int {
	return int(endTime.Sub(startTime).Minutes())
}

func toAverageDeliveryTimeSlice(items []movingAverage) []domain.AverageDeliveryTime {
	domainItems := make([]domain.AverageDeliveryTime, len(items))

	for i := 0; i < len(items); i++ {
		item := items[i]
		domainItems[i] = domain.AverageDeliveryTime{
			Date:                item.Timestamp,
			AverageDeliveryTime: item.AverageDuration,
		}
	}

	return domainItems
}
