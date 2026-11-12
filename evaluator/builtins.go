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
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},

	// Common utilities
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

	// Strings standard library
	"str_to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for str_to_upper. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("str_to_upper expects string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for str_to_lower. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("str_to_lower expects string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for str_contains. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("str_contains expects two strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for str_split. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("str_split expects two strings")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for str_trim. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("str_trim expects string")
			}
			return &object.String{Value: strings.TrimSpace(s.Value)}
		},
	},

	// Math standard library
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_sqrt expects 1 argument")
			}
			val, ok := getNumberAsFloat(args[0])
			if !ok {
				return newError("math_sqrt expects a numeric argument")
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_pow expects 2 arguments")
			}
			base, ok1 := getNumberAsFloat(args[0])
			exponent, ok2 := getNumberAsFloat(args[1])
			if !ok1 || !ok2 {
				return newError("math_pow expects numeric arguments")
			}
			return &object.Float{Value: math.Pow(base, exponent)}
		},
	},
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_abs expects 1 argument")
			}
			switch val := args[0].(type) {
			case *object.Integer:
				v := val.Value
				if v < 0 {
					v = -v
				}
				return &object.Integer{Value: v}
			case *object.Float:
				return &object.Float{Value: math.Abs(val.Value)}
			default:
				return newError("math_abs expects a numeric argument")
			}
		},
	},
	"math_floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_floor expects 1 argument")
			}
			val, ok := getNumberAsFloat(args[0])
			if !ok {
				return newError("math_floor expects a numeric argument")
			}
			return &object.Integer{Value: int64(math.Floor(val))}
		},
	},
	"math_ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_ceil expects 1 argument")
			}
			val, ok := getNumberAsFloat(args[0])
			if !ok {
				return newError("math_ceil expects a numeric argument")
			}
			return &object.Integer{Value: int64(math.Ceil(val))}
		},
	},
	"math_min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_min expects 2 arguments")
			}
			v1, ok1 := getNumberAsFloat(args[0])
			v2, ok2 := getNumberAsFloat(args[1])
			if !ok1 || !ok2 {
				return newError("math_min expects numeric arguments")
			}
			if _, isInt1 := args[0].(*object.Integer); isInt1 {
				if _, isInt2 := args[1].(*object.Integer); isInt2 {
					if v1 < v2 {
						return args[0]
					}
					return args[1]
				}
			}
			return &object.Float{Value: math.Min(v1, v2)}
		},
	},
	"math_max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_max expects 2 arguments")
			}
			v1, ok1 := getNumberAsFloat(args[0])
			v2, ok2 := getNumberAsFloat(args[1])
			if !ok1 || !ok2 {
				return newError("math_max expects numeric arguments")
			}
			if _, isInt1 := args[0].(*object.Integer); isInt1 {
				if _, isInt2 := args[1].(*object.Integer); isInt2 {
					if v1 > v2 {
						return args[0]
					}
					return args[1]
				}
			}
			return &object.Float{Value: math.Max(v1, v2)}
		},
	},

	// Time standard library
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixMilli()}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("time_sleep expects 1 argument (milliseconds)")
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("time_sleep expects integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("read_file expects 1 argument (filepath)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("read_file expects a string filepath")
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
				return newError("write_file expects 2 arguments (filepath, content)")
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("write_file expects string filepath and string content")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file '%s': %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func getNumberAsFloat(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
