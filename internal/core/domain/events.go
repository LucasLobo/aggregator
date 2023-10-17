package domain

// TranslationDelivered is an event that represents the delivery time of a translation.
type TranslationDelivered struct {
	Timestamp      Time   `json:"timestamp"`
	TranslationId  string `json:"translation_id"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	ClientName     string `json:"client_name"`
	EventName      string `json:"event_name"`
	NrWords        int    `json:"nr_words"`
	Duration       int    `json:"duration"`
}

// AverageDeliveryTime represents the average delivery time.
type AverageDeliveryTime struct {
	Date                Time    `json:"date"`
	AverageDeliveryTime float32 `json:"average_delivery_time"`
}
