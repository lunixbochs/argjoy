package argjoy

import (
	"fmt"
	"strconv"
	"strings"
)

func strBaseToInt(arg interface{}, vals []interface{}, base int) error {
	if s, ok := vals[0].(string); ok {
		// we only need to ParseInt on a leading negative sign
		if len(s) > 0 && s[0] == '-' {
			// parsing signed int
			n, err := strconv.ParseInt(s, base, 64)
			if err != nil {
				return err
			}
			return IntToInt(arg, []interface{}{n})
		} else if base == 0 && strings.HasPrefix(s, "0b") {
			// parse base-2 int
			n, err := strconv.ParseUint(s[2:], 2, 64)
			if err != nil {
				return err
			}
			return IntToInt(arg, []interface{}{n})
		} else {
			// parsing unsigned int
			n, err := strconv.ParseUint(s, base, 64)
			if err != nil {
				return err
			}
			return IntToInt(arg, []interface{}{n})
		}
	}
	return NoMatch
}

// Codec to convert a string to any int type (base 10 forced).
// Throws bounds and sign errors as appropriate.
func StrToInt(arg interface{}, vals []interface{}) error {
	return strBaseToInt(arg, vals, 10)
}

// Codec to convert radix notations like hex, octal, binary to int.
// Throws bounds and sign errors as appropriate.
func RadStrToInt(arg interface{}, vals []interface{}) error {
	return strBaseToInt(arg, vals, 0)
}

// Codec to convert between any builtin int types.
// Throws bounds and sign errors as appropriate.
func IntToInt(arg interface{}, vals []interface{}) (err error) {
	var ival int64
	var uval uint64
	switch src := vals[0].(type) {
	case int:
		ival = int64(src)
	case int8:
		ival = int64(src)
	case int16:
		ival = int64(src)
	case int32:
		ival = int64(src)
	case int64:
		ival = int64(src)
	case uint:
		uval = uint64(src)
	case uint8:
		uval = uint64(src)
	case uint16:
		uval = uint64(src)
	case uint32:
		uval = uint64(src)
	case uint64:
		uval = uint64(src)
	default:
		return NoMatch
	}

	// prep bounds check
	var (
		min int64
		max uint64
	)
	switch arg.(type) {
	case *int:
		// 32-bit vs 64-bit
		if ^uint(0) == (1<<32)-1 {
			max = (1 << 31) - 1
		} else {
			max = (1 << 63) - 1
		}
		min = int64(^max)
	case *int8:
		max = (1 << 7) - 1
		min = int64(^max)
	case *int16:
		max = (1 << 15) - 1
		min = int64(^max)
	case *int32:
		max = (1 << 31) - 1
		min = int64(^max)
	case *int64:
		max = (1 << 63) - 1
		min = int64(^max)
	case *uint:
		max = uint64(^uint(0))
	case *uint8:
		max = uint64(^uint8(0))
	case *uint16:
		max = uint64(^uint16(0))
	case *uint32:
		max = uint64(^uint32(0))
	case *uint64:
		max = ^uint64(0)
	default:
		return NoMatch
	}

	// assignment
	if ival != 0 {
		// signed
		// the sign bit is cleared on max to prevent the cast from forcing it negative
		if ival < min || ival > int64(max&^(1<<63)) {
			return fmt.Errorf("int arg (%d) outside range for type (%T)", ival, arg)
		}
		switch a := arg.(type) {
		case *int:
			*a = int(ival)
		case *int8:
			*a = int8(ival)
		case *int16:
			*a = int16(ival)
		case *int32:
			*a = int32(ival)
		case *int64:
			*a = int64(ival)
		case *uint:
			*a = uint(ival)
		case *uint8:
			*a = uint8(ival)
		case *uint16:
			*a = uint16(ival)
		case *uint32:
			*a = uint32(ival)
		case *uint64:
			*a = uint64(ival)
		}
	} else {
		if uval > max {
			return fmt.Errorf("uint arg (%d) outside range for type (%T)", uval, arg)
		}
		// unsigned
		switch a := arg.(type) {
		case *int:
			*a = int(uval)
		case *int8:
			*a = int8(uval)
		case *int16:
			*a = int16(uval)
		case *int32:
			*a = int32(uval)
		case *int64:
			*a = int64(uval)
		case *uint:
			*a = uint(uval)
		case *uint8:
			*a = uint8(uval)
		case *uint16:
			*a = uint16(uval)
		case *uint32:
			*a = uint32(uval)
		case *uint64:
			*a = uint64(uval)
		}
	}
	return nil
}
