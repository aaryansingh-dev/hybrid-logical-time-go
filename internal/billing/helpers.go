package billing

import (
	"reflect"
)

// getEventName retrieves the string representation of a struct's type name.
// It handles both value and pointer types using reflection, allowing the
// engine to generate human-readable logs without manual string mapping.
func getEventName(e any) string {
	t := reflect.TypeOf(e)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
