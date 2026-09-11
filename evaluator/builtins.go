package evaluator

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

func RegisterBuiltinsWithWriter(out io.Writer) map[string]*object.Builtin {
	b := map[string]*object.Builtin{
		// Printing
		"println": {
			Fn: func(args ...object.Object) object.Object {
				var parts []string
				for _, arg := range args {
					parts = append(parts, arg.Inspect())
				}
				fmt.Fprintln(out, strings.Join(parts, " "))
				return &object.Null{}
			},
		},
		"print": {
			Fn: func(args ...object.Object) object.Object {
				var parts []string
				for _, arg := range args {
					parts = append(parts, arg.Inspect())
				}
				fmt.Fprint(out, strings.Join(parts, " "))
				return &object.Null{}
			},
		},

		// Strings / Collection Builtins
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
		"contains": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				str, ok1 := args[0].(*object.String)
				substr, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("arguments to `contains` must be strings")
				}
				return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
			},
		},
		"split": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				str, ok1 := args[0].(*object.String)
				sep, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("arguments to `split` must be strings")
				}
				parts := strings.Split(str.Value, sep.Value)
				elements := make([]object.Object, len(parts))
				for i, p := range parts {
					elements[i] = &object.String{Value: p}
				}
				return &object.Array{Elements: elements}
			},
		},
		"join": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				arr, ok1 := args[0].(*object.Array)
				sep, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("first argument to `join` must be an array and second a string")
				}
				var parts []string
				for _, elem := range arr.Elements {
					parts = append(parts, elem.Inspect())
				}
				return &object.String{Value: strings.Join(parts, sep.Value)}
			},
		},
		"to_upper": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				str, ok := args[0].(*object.String)
				if !ok {
					return newError("argument to `to_upper` must be string")
				}
				return &object.String{Value: strings.ToUpper(str.Value)}
			},
		},
		"to_lower": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				str, ok := args[0].(*object.String)
				if !ok {
					return newError("argument to `to_lower` must be string")
				}
				return &object.String{Value: strings.ToLower(str.Value)}
			},
		},

		// Math Builtins
		"abs": {
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
					return newError("argument to `abs` must be integer or float")
				}
			},
		},
		"pow": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				base, ok1 := getFloatValue(args[0])
				exponent, ok2 := getFloatValue(args[1])
				if !ok1 || !ok2 {
					return newError("arguments to `pow` must be numbers")
				}
				res := math.Pow(base, exponent)
				if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
					return &object.Integer{Value: int64(res)}
				}
				return &object.Float{Value: res}
			},
		},
		"sin": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				v, ok := getFloatValue(args[0])
				if !ok {
					return newError("argument to `sin` must be a number")
				}
				return &object.Float{Value: math.Sin(v)}
			},
		},
		"cos": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				v, ok := getFloatValue(args[0])
				if !ok {
					return newError("argument to `cos` must be a number")
				}
				return &object.Float{Value: math.Cos(v)}
			},
		},
		"tan": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				v, ok := getFloatValue(args[0])
				if !ok {
					return newError("argument to `tan` must be a number")
				}
				return &object.Float{Value: math.Tan(v)}
			},
		},

		// Time Builtins
		"now": {
			Fn: func(args ...object.Object) object.Object {
				return &object.Integer{Value: time.Now().UnixNano()}
			},
		},
		"sleep": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return newError("argument to `sleep` must be integer (milliseconds)")
				}
				time.Sleep(time.Duration(ms.Value) * time.Millisecond)
				return &object.Null{}
			},
		},

		// Filesystem Basics
		"read_file": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return newError("argument to `read_file` must be string")
				}
				content, err := os.ReadFile(path.Value)
				if err != nil {
					return newError("could not read file %s: %s", path.Value, err.Error())
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
					return newError("arguments to `write_file` must be strings")
				}
				err := os.WriteFile(path.Value, []byte(content.Value), 0644)
				if err != nil {
					return newError("could not write file %s: %s", path.Value, err.Error())
				}
				return &object.Boolean{Value: true}
			},
		},
	}
	return b
}

var builtins = RegisterBuiltinsWithWriter(os.Stdout)

func getFloatValue(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
