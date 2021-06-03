package internal

import (
	"encoding/json"
)

func Parse(input string) Event {
	var ie Event
	json.Unmarshal([]byte(input), &ie)

	return ie
}
