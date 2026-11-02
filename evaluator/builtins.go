package evaluator

import (
	"fmt"
	"math"
	"os"
	"rustx/object"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	// Printing
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	// Strings standard library
	"str_len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `str_len` must be STRING, got %s", args[0].Type())
			}
			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `str_contains` must be (STRING, STRING)")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `str_split` must be (STRING, STRING)")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `str_upper` must be STRING")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `str_lower` must be STRING")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `str_trim` must be STRING")
			}
			return &object.String{Value: strings.TrimSpace(s.Value)}
		},
	},

	// Math standard library
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch a := args[0].(type) {
			case *object.Integer:
				val = float64(a.Value)
			case *object.Float:
				val = a.Value
			default:
				return newError("argument to `math_sqrt` must be number, got %s", args[0].Type())
			}
			if val < 0 {
				return newError("cannot take square root of negative number: %g", val)
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch a := args[0].(type) {
			case *object.Integer:
				if a.Value < 0 {
					return &object.Integer{Value: -a.Value}
				}
				return a
			case *object.Float:
				return &object.Float{Value: math.Abs(a.Value)}
			default:
				return newError("argument to `math_abs` must be number, got %s", args[0].Type())
			}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, exp float64
			switch a := args[0].(type) {
			case *object.Integer:
				base = float64(a.Value)
			case *object.Float:
				base = a.Value
			default:
				return newError("first argument to `math_pow` must be number")
			}
			switch a := args[1].(type) {
			case *object.Integer:
				exp = float64(a.Value)
			case *object.Float:
				exp = a.Value
			default:
				return newError("second argument to `math_pow` must be number")
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"math_min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			v1, ok1 := toFloat(args[0])
			v2, ok2 := toFloat(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `math_min` must be numbers")
			}
			if v1 < v2 {
				return args[0]
			}
			return args[1]
		},
	},
	"math_max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			v1, ok1 := toFloat(args[0])
			v2, ok2 := toFloat(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `math_max` must be numbers")
			}
			if v1 > v2 {
				return args[0]
			}
			return args[1]
		},
	},

	// Time standard library
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			dur, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `time_sleep` must be INTEGER milliseconds")
			}
			time.Sleep(time.Duration(dur.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be STRING filename")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("failed to read file '%s': %s", path.Value, err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be (STRING filename, STRING content)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file '%s': %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func toFloat(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
