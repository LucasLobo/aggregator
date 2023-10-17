package domain

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	inputTimeLayout  = "2006-01-02 15:04:05.999999"
	outputTimeLayout = "2006-01-02 15:04:00"
)

// Time is a custom time type that allows us to marshall and unmarshall with the specific formats expected
// in the input and output
type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(outputTimeLayout))
}

func (t *Time) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsed, err := time.Parse(inputTimeLayout, str)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}
