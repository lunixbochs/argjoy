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

// Any codecs passed to NewArgjoy will be passed to Register() on the new instance.
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

func (a *Argjoy) translate(arg, val interface{}) error {
	matched := false
	// O(N) :( only way to avoid this is an arg type registry?
	for _, codec := range a.codecs {
		err := codec(arg, val)
		if err == nil {
			matched = true
			break
		} else if _, ok := err.(*NoMatchErr); !ok {
			return err
		}
	}
	if !matched {
		return &NoMatchErr{fmt.Sprintf("%T -> %T", val, arg)}
	}
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
		if i == argCount-1 && fnt.IsVariadic() {
			// this condition is nested because we don't want to ever fall to the next
			// else block on variadic functions
			if len(vals) >= argCount {
				varargs := reflect.MakeSlice(argType, len(vals[i:]), len(vals[i:]))
				arg.Elem().Set(varargs)
				for j, v := range vals[i:] {
					// we need this because New makes a pointer
					arg := varargs.Index(j).Addr()
					val := reflect.ValueOf(v)
					if err := a.translate(arg.Interface(), val.Interface()); err != nil {
						return nil, err
					}
				}
			}
		} else if i < len(vals) {
			val := reflect.ValueOf(vals[i])
			if err := a.translate(arg.Interface(), val.Interface()); err != nil {
				return nil, err
			}
		} else if a.Optional {
			// we ran out of input and Optional arguments are enabled
			// so the rest of the args are zeroed (which happens automatically on reflect.New())
		} else {
			// we ran out of input and Optional arguments are disabled
			return nil, &ArgCountErr{len(vals), argCount}
		}
		in[i] = arg.Elem()
	}
	var out []reflect.Value
	if fnt.IsVariadic() {
		out = fnv.CallSlice(in)
	} else {
		out = fnv.Call(in)
	}
	ret := make([]interface{}, len(out))
	for i, v := range out {
		ret[i] = v.Interface()
	}
	return ret, nil
}
