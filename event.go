package events

import (
	"encoding/json"
	"errors"
	"reflect"
)

// Event represents different events
// in the lifecycle of a Buffalo app
type Event struct {
	// Kind is the "type" of event "app:start"
	Kind string `json:"kind"`
	// Message is optional
	Message string `json:"message"`
	// Payload is optional
	Payload interface{} `json:"payload"`
	// Error is optional
	Error error `json:"-"`
}

func (e Event) String() string {
	b, _ := e.MarshalJSON()

	return string(b)
}

// MarshalJSON implements the json marshaler for an event
func (e Event) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"kind": e.Kind,
	}
	if len(e.Message) != 0 {
		m["message"] = e.Message
	}
	if e.Error != nil {
		m["error"] = e.Error.Error()
	}

	rv := reflect.Indirect(reflect.ValueOf(e.Payload))
	switch rv.Kind() {
	case reflect.Map:
		pm := map[string]interface{}{}
		for _, k := range rv.MapKeys() {
			v := rv.MapIndex(k)
			if _, err := json.Marshal(v.Interface()); err == nil {
				// it can be marshaled, so add it:
				pm[k.String()] = v.Interface()
			}
		}
		m["payload"] = pm
	default:
		if _, err := json.Marshal(e.Payload); err == nil {
			// it can be marshaled, so add it:
			m["payload"] = e.Payload
		}
	}

	return json.Marshal(m)
}

// Validate that an event is ready to be emitted
func (e Event) Validate() error {
	if len(e.Kind) == 0 {
		return errors.New("kind can not be blank")
	}
	return nil
}
