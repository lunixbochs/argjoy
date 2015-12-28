package argjoy

import (
	"fmt"
	"reflect"
)

type CodecFunc func(arg interface{}, vals []interface{}) error

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

func (a *Argjoy) translate(arg interface{}, vals []interface{}) error {
	matched := false
	// O(N) :( only way to avoid this is an arg type registry?
	for _, codec := range a.codecs {
		err := codec(arg, vals)
		if err == nil {
			matched = true
			break
		} else if _, ok := err.(*NoMatchErr); !ok {
			return err
		}
	}
	if !matched {
		// if types are equal, at this point assign arg = vals[0]
		// TODO: maybe make this step configurable
		if reflect.TypeOf(arg).Elem() == reflect.TypeOf(vals[0]) {
			reflect.ValueOf(arg).Elem().Set(reflect.ValueOf(vals[0]))
			return nil
		}
		return &NoMatchErr{fmt.Sprintf("%T -> %T", vals[0], arg)}
	}
	return nil
}

// Given the arguments to a function as per reflect.In(),
// produce a list of arguments suitable for invoking reflect.Call(fn, args).
// If fn is variadic, you must instead use reflect.CallSlice(fn, args).
func (a *Argjoy) Convert(in []reflect.Type, variadic bool, vals ...interface{}) ([]reflect.Value, error) {
	var vals2 []interface{}
	for _, v := range vals {
		slice := reflect.ValueOf(v)
		if slice.Kind() == reflect.Slice {
			for i := 0; i < slice.Len(); i++ {
				vals2 = append(vals2, slice.Index(i).Interface())
			}
		} else {
			vals2 = append(vals2, v)
		}
	}
	vals = vals2

	ret := make([]reflect.Value, len(in))
	for i, argType := range in {
		arg := reflect.New(argType)
		if i == len(in)-1 && variadic {
			// this condition is nested because we don't want to ever fall to the next
			// else block on variadic functions
			if len(vals) >= len(in) {
				varargs := reflect.MakeSlice(argType, len(vals[i:]), len(vals[i:]))
				arg.Elem().Set(varargs)
				for j := range vals[i:] {
					// we need this because New makes a pointer
					arg := varargs.Index(j).Addr()
					if err := a.translate(arg.Interface(), vals[i+j:]); err != nil {
						return nil, err
					}
				}
			}
		} else if i < len(vals) {
			if err := a.translate(arg.Interface(), vals[i:]); err != nil {
				return nil, err
			}
		} else if a.Optional {
			// we ran out of input and Optional arguments are enabled
			// so the rest of the args are zeroed (which happens automatically on reflect.New())
		} else {
			// we ran out of input and Optional arguments are disabled
			return nil, &ArgCountErr{len(vals), len(in)}
		}
		ret[i] = arg.Elem()
	}
	return ret, nil
}

// Call fn(vals...), using registered codecs to convert argument types.
// Returns []interface{} of target function's return values.
// Will return an error value if any codec fails.
// Panics if fn is not a valid function.
func (a *Argjoy) Call(fn interface{}, vals ...interface{}) ([]interface{}, error) {
	fnv := reflect.ValueOf(fn)
	fnt := fnv.Type()
	if fnv.Kind() != reflect.Func {
		panic(fmt.Sprintf("function expected as first argument, got: (%T)%#v", fn, fn))
	}
	inTypes := make([]reflect.Type, fnt.NumIn())
	for i := range inTypes {
		inTypes[i] = fnt.In(i)
	}
	in, err := a.Convert(inTypes, fnt.IsVariadic(), vals...)
	if err != nil {
		return nil, err
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
