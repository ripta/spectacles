package sinks

import (
	"encoding/json"

	apiv1 "k8s.io/api/core/v1"
)

// Encoder is a function that can encode an event into a sequence of bytes.
type Encoder func(*apiv1.Event) ([]byte, error)

// JSONEncoder encodes an event into JSON.
func JSONEncoder(e *apiv1.Event) ([]byte, error) {
	return json.Marshal(e)
}
