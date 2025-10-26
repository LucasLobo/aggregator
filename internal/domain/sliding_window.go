package domain

import (
	"time"
)

type SlidingWindow struct {
	windowSize int
	buckets    map[time.Time]state
	state      state
	start      time.Time
	head       time.Time
	tail       time.Time
}

type state struct {
	count    int
	duration int
}

func NewSlidingWindow(windowSize int) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		buckets:    map[time.Time]state{},
	}

}

func (sw *SlidingWindow) Process(event DurationEvent) []AverageDeliveryTime {

	bucket := event.Timestamp.Truncate(time.Minute).Add(time.Minute)

	if sw.start.IsZero() {
		start := bucket.Add(-time.Minute)
		sw.start = start
		sw.head = start
		sw.tail = start
	}

	var adt []AverageDeliveryTime
	for beforeOrEqual(sw.head, bucket) {
		sw.buckets[sw.head] = state{}
		// when we are at the time bucket of the current event, we add it to the state
		if sw.head.Equal(bucket) {
			sw.state.count += 1
			sw.state.duration += event.Duration
			sw.buckets[sw.head] = state{
				duration: event.Duration,
				count:    1,
			}
		}

		// when we exceed the current window size, we must remove the last item and advance the tail
		if len(sw.buckets) > sw.windowSize {
			sw.state.count -= sw.buckets[sw.tail].count
			sw.state.duration -= sw.buckets[sw.tail].duration
			delete(sw.buckets, sw.tail)
			sw.tail = sw.tail.Add(time.Minute)
		}

		// once we're done, we calculate the average for the current position
		average := float32(0)
		if sw.state.count != 0 {
			// let's not divide by 0 ;)
			average = float32(sw.state.duration) / float32(sw.state.count)

		}
		adt = append(adt, AverageDeliveryTime{
			Date:                Time{Time: sw.head},
			AverageDeliveryTime: average,
		})

		// at the end we must advance the head to keep going
		sw.head = sw.head.Add(time.Minute)
	}
	return adt

}

func beforeOrEqual(a, b time.Time) bool {
	return a.Before(b) || a.Equal(b)
}
