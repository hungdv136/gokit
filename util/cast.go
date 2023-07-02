package util

import (
	"fmt"
)

// ToArrayString converts from array of interface{} to array of string
func ToArrayString(a []interface{}) ([]string, bool) {
	var ok bool
	b := make([]string, len(a))
	for i, e := range a {
		b[i], ok = e.(string)
		if !ok {
			return nil, false
		}
	}
	return b, true
}

// ToString converts anything to string
func ToString(any interface{}) string {
	switch any.(type) {
	case float32, float64:
		return fmt.Sprintf("%.6f", any)
	default:
		return fmt.Sprintf("%v", any)
	}
}
