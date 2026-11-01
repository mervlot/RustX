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
			var parts []string
			for _, arg := range args {
				if str, ok := arg.(*object.String); ok {
					parts = append(parts, str.Value)
				} else {
					parts = append(parts, arg.Inspect())
				}
			}
			fmt.Println(strings.Join(parts, " "))
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				if str, ok := arg.(*object.String); ok {
					parts = append(parts, str.Value)
				} else {
					parts = append(parts, arg.Inspect())
				}
			}
			fmt.Print(strings.Join(parts, " "))
			return NULL
		},
	},

	// Builtin collections
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1])
			return arr
		},
	},

	// Strings module functions
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("str_contains requires 2 arguments (s, substr)")
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_contains must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	"str_to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str_to_upper requires 1 string argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_to_upper must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str_to_lower requires 1 string argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_to_lower must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("str_split requires 2 arguments (s, sep)")
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_split must be strings")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},

	// Math module functions
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_sqrt requires 1 number argument")
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("math_sqrt requires numeric argument, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_abs requires 1 number argument")
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value < 0 {
					return &object.Integer{Value: -arg.Value}
				}
				return arg
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError("math_abs requires numeric argument")
			}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_pow requires 2 number arguments (base, exp)")
			}
			var b, exp float64
			switch arg := args[0].(type) {
			case *object.Integer:
				b = float64(arg.Value)
			case *object.Float:
				b = arg.Value
			default:
				return newError("math_pow base must be numeric")
			}
			switch arg := args[1].(type) {
			case *object.Integer:
				exp = float64(arg.Value)
			case *object.Float:
				exp = arg.Value
			default:
				return newError("math_pow exp must be numeric")
			}
			return &object.Float{Value: math.Pow(b, exp)}
		},
	},

	// Time module functions
	"time_now_unix": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_sleep_ms": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("time_sleep_ms requires 1 argument (duration in ms)")
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to time_sleep_ms must be integer")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem module functions
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("read_file requires 1 string argument (filepath)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("filepath must be a string")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("read_file failed: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("write_file requires 2 arguments (filepath, content)")
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("write_file arguments must be strings")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("write_file failed: %s", err.Error())
			}
			return TRUE
		},
	},
}
