package evaluator

import (
	"fmt"
	"math"
	"os"
	"rustx/object"
	"strconv"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	// Printing
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
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

	// Strings module functions
	"str_len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_len must be STRING, got %s", args[0].Type())
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
				return newError("arguments to str_contains must be STRING")
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
				return newError("arguments to str_split must be STRING")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_to_upper must be STRING")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_to_lower must be STRING")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"to_string": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: args[0].Inspect()}
		},
	},
	"to_int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				return arg
			case *object.Float:
				return &object.Integer{Value: int64(arg.Value)}
			case *object.String:
				v, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("cannot parse %q as int", arg.Value)
				}
				return &object.Integer{Value: v}
			default:
				return newError("cannot convert %s to int", arg.Type())
			}
		},
	},
	"to_float": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Float:
				return arg
			case *object.Integer:
				return &object.Float{Value: float64(arg.Value)}
			case *object.String:
				v, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("cannot parse %q as float", arg.Value)
				}
				return &object.Float{Value: v}
			default:
				return newError("cannot convert %s to float", arg.Type())
			}
		},
	},

	// Math library functions
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Float:
				val = arg.Value
			case *object.Integer:
				val = float64(arg.Value)
			default:
				return newError("argument to math_sqrt must be number, got %s", args[0].Type())
			}
			if val < 0 {
				return newError("cannot take square root of negative number: %g", val)
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, exp float64
			switch b := args[0].(type) {
			case *object.Float:
				base = b.Value
			case *object.Integer:
				base = float64(b.Value)
			default:
				return newError("base to math_pow must be number, got %s", args[0].Type())
			}
			switch e := args[1].(type) {
			case *object.Float:
				exp = e.Value
			case *object.Integer:
				exp = float64(e.Value)
			default:
				return newError("exponent to math_pow must be number, got %s", args[1].Type())
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
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
				return newError("argument to math_abs must be number, got %s", args[0].Type())
			}
		},
	},
	"math_max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			v1, ok1 := args[0].(*object.Integer)
			v2, ok2 := args[1].(*object.Integer)
			if ok1 && ok2 {
				if v1.Value >= v2.Value {
					return v1
				}
				return v2
			}
			return newError("arguments to math_max must currently be integers")
		},
	},
	"math_min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			v1, ok1 := args[0].(*object.Integer)
			v2, ok2 := args[1].(*object.Integer)
			if ok1 && ok2 {
				if v1.Value <= v2.Value {
					return v1
				}
				return v2
			}
			return newError("arguments to math_min must currently be integers")
		},
	},

	// Time library functions
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_now_millis": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			millis, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to time_sleep must be INTEGER milliseconds")
			}
			time.Sleep(time.Duration(millis.Value) * time.Millisecond)
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
				return newError("argument to read_file must be STRING")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("could not read file %q: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
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
				return newError("arguments to write_file must be STRING")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("could not write file %q: %s", path.Value, err.Error())
			}
			return nativeBoolToBooleanObject(true)
		},
	},
}
