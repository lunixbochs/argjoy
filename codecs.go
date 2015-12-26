package argjoy

import (
	"strconv"
)

// Codec to convert a string to any int type.
// Throws bounds and sign errors as appropriate
func StrToInt(arg, val interface{}) (err error) {
	if s, ok := val.(string); ok {
		switch a := arg.(type) {
		case *uint:
			var n uint64
			// we don't actually know uint and int sizes without doing math
			if ^(uint(0)) == (1<<32)-1 {
				n, err = strconv.ParseUint(s, 10, 32)
			} else {
				n, err = strconv.ParseUint(s, 10, 64)
			}
			*a = uint(n)
		case *int:
			var n int64
			if ^(uint(0)) == (1<<32)-1 {
				n, err = strconv.ParseInt(s, 10, 32)
			} else {
				n, err = strconv.ParseInt(s, 10, 64)
			}
			*a = int(n)
		case *uint32:
			var n uint64
			n, err = strconv.ParseUint(s, 10, 32)
			*a = uint32(n)
		case *uint64:
			*a, err = strconv.ParseUint(s, 10, 64)
		case *int32:
			var n int64
			n, err = strconv.ParseInt(s, 10, 32)
			*a = int32(n)
		case *int64:
			*a, err = strconv.ParseInt(s, 10, 64)
		default:
			return NoMatchErr
		}
		return
	}
	return NoMatchErr
}
