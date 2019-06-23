package events

import (
	"regexp"

	"github.com/gobuffalo/events/internal/safe"
)

// Filter compiles the string as a regex and returns
// the original listener wrapped in a new listener
// that filters incoming events by the Kind
func Filter(s string, fn Listener) Listener {
	if s == "" || s == "*" {
		return fn
	}
	rx := regexp.MustCompile(s)
	return func(e Event) {
		if rx.MatchString(e.Kind) {
			safe.Run(func() {
				fn(e)
			})
		}
	}
}
