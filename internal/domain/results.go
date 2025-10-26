package domain

// AverageDeliveryTime represents the average delivery time for the sliding window
type AverageDeliveryTime struct {
	Date                Time    `json:"date"`
	AverageDeliveryTime float32 `json:"average_delivery_time"`
}
