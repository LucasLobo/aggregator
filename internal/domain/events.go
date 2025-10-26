package domain

// DurationEvent is an event that represents the delivery time of a translation.
type DurationEvent struct {
	Timestamp Time `json:"timestamp"`
	Duration  int  `json:"duration"`
}
