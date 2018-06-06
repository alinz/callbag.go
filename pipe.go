package callbag

import (
	"reflect"
)

// Pipe is a utility function for plugging callbags together in chain.
//
func Pipe(cbs ...interface{}) Source {
	var res Source

	for i := 0; i < len(cbs); i++ {
		cb := cbs[i]
		if cb == nil || reflect.ValueOf(cb).IsNil() {
			panic("cb is nil")
		}

		switch v := cb.(type) {
		case Source:
			if i != 0 {
				panic("Source must be the first argument")
			}
			res = v

		case Transform:
			res = v(res)

		case Sink:
			if i != len(cbs)-1 {
				panic("Sink must be the last argument")
			}
			v(res)
			res = nil

		default:
			panic("unknown type")
		}
	}

	return res
}
