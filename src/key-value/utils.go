package key_value

import (
	"fmt"
	"reflect"
)

func Equal(a, b any) bool {
	t1 := reflect.TypeOf(a)
	t2 := reflect.TypeOf(b)

	if t1 != t2 {
		return false
	}

	switch a.(type) {
	case nil:
		return true
	case bool:
		return a == b
	case float64:
		return a == b
	case string:
		return a == b
	case []any:
		aa := a.([]any)
		bb := b.([]any)

		if len(aa) != len(bb) {
			return false
		}

		for i, v := range aa {
			if !equal(v, bb[i]) {
				return false
			}
		}

		return true
	case map[string]any:
		aa := a.(map[string]any)
		bb := b.(map[string]any)

		if len(aa) != len(bb) {
			return false
		}

		for k, v := range aa {
			if !equal(v, bb[k]) {
				return false
			}
		}

		return true
	default:
		panic(fmt.Errorf("unknown type of values to compare: %T-%T", a, b))
	}
}
