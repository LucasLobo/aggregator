package domain

import (
	"encoding/json"
	"time"
)

// TranslationDelivered is an event that represents the delivery time of a translation.
type TranslationDelivered struct {
	Timestamp      time.Time `json:"timestamp"`
	TranslationId  string    `json:"translation_id"`
	SourceLanguage string    `json:"source_language"`
	TargetLanguage string    `json:"target_language"`
	ClientName     string    `json:"client_name"`
	EventName      string    `json:"event_name"`
	NrWords        int       `json:"nr_words"`
	Duration       int       `json:"duration"`
}

// AverageDeliveryTime represents the average delivery time.
type AverageDeliveryTime struct {
	Date                time.Time `json:"date"`
	AverageDeliveryTime float32   `json:"average_delivery_time"`
}

// We define Marshalling and Unmarshalling methods to use the specific date formats expects in the events

func (t *TranslationDelivered) UnmarshalJSON(data []byte) error {
	type Alias TranslationDelivered
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Timestamp, _ = time.Parse("2006-01-02 15:04:05.999999", aux.Timestamp)
	return nil
}

func (t AverageDeliveryTime) MarshalJSON() ([]byte, error) {
	type Alias AverageDeliveryTime
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  t.Date.Format("2006-01-02 15:04:00"),
		Alias: (*Alias)(&t),
	})
}
