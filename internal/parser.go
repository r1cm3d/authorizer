package internal

import (
	"encoding/json"
	"strings"
	"time"
)

func Parse(input string) Event {
	var ie Event
	json.Unmarshal([]byte(input), &ie)

	return ie
}

func (it *Time) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, _ := time.Parse(time.RFC3339, s)

	*it = Time(t)
	return nil
}
