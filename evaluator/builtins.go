package evaluator

import (
	"fmt"
	"math"

	"github.com/username/rustx/object"
)

var builtins = map[string]*object.Builtin{
	"println": {
		Fn: func(args ...object.Object) object.Object {
			out := []interface{}{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(out...)
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			out := []interface{}{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(out...)
			return NULL
		},
	},
	// Strings builtins
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
	// Math builtins
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Integer{Value: val}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			base, ok1 := getFloatValue(args[0])
			exp, ok2 := getFloatValue(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `pow` must be numbers")
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to `sin` must be a number")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to `cos` must be a number")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to `tan` must be a number")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
}

func getFloatValue(obj object.Object) (float64, bool) {
	switch val := obj.(type) {
	case *object.Integer:
		return float64(val.Value), true
	case *object.Float:
		return val.Value, true
	default:
		return 0, false
	}
}
