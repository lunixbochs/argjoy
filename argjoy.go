package argjoy

import (
	"fmt"
	"reflect"
)

type CodecFunc func(arg, val interface{}) error

type Argjoy struct {
	codecs []CodecFunc

	// Determines whether unpassed arguments are optional.
	// If this is false, invoking Call() with insufficient number of arguments
	// will return ArgCountErr
	Optional bool
}

func NewArgjoy(codecs ...CodecFunc) *Argjoy {
	n := &Argjoy{}
	for _, codec := range codecs {
		n.Register(codec)
	}
	return n
}

// Registers a new codec function which will be used to convert arguments during Call()
func (a *Argjoy) Register(codec CodecFunc) error {
	a.codecs = append(a.codecs, codec)
	return nil
}

// Call fn(vals...), using registered codecs to convert argument types.
// Returns []interface{} of target function's return values.
// Will return an error value if any codec fails.
// Panics if fn is not a valid function.
func (a *Argjoy) Call(fn interface{}, vals ...interface{}) ([]interface{}, error) {
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		panic(fmt.Sprintf("Argjoy.Call() requires function as first argument, got: (%T)%#v", fn, fn))
	}
	fnt := fnv.Type()
	argCount := fnt.NumIn()
	in := make([]reflect.Value, argCount)
	for i := 0; i < argCount; i++ {
		argType := fnt.In(i)
		arg := reflect.New(argType)
		if i < len(vals) {
			val := reflect.ValueOf(vals[i])
			matched := false
			// O(N) :( only way to avoid this is an arg type registry?
			for _, codec := range a.codecs {
				err := codec(arg.Interface(), val.Interface())
				if err == nil {
					matched = true
					break
				} else if err != NoMatchErr {
					return nil, err
				}
			}
			if !matched {
				return nil, NoMatchErr
			}
		} else if a.Optional {
			// we ran out of input and Optional arguments are enabled
			// so the rest of the args are zeroed (which happens automatically on reflect.New())
		} else {
			// we ran out of input and Optional arguments are disabled
			// so... panic or error?
			return nil, ArgCountErr
		}
		in[i] = arg.Elem()
	}
	out := fnv.Call(in)
	ret := make([]interface{}, len(out))
	for i, v := range out {
		ret[i] = v.Interface()
	}
	return ret, nil
}
