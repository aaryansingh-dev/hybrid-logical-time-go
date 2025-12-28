package billing

import (
	"reflect"
)

// GetTypeName returns the type name of the struct
func getEventName(e any) string {
	t := reflect.TypeOf(e)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
