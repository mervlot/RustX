package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

func registerBuiltins() {
	// print
	builtins["print"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Print(strings.Join(parts, " "))
			return NULL
		},
	}

	// println
	builtins["println"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(strings.Join(parts, " "))
			return NULL
		},
	}

	// len
	builtins["len"] = &object.Builtin{
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
	}

	// abs
	builtins["abs"] = &object.Builtin{
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
				return newError("argument to `abs` must be integer or float, got %s", args[0].Type())
			}
		},
	}

	// pow
	builtins["pow"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			x, okX := args[0].(*object.Float)
			y, okY := args[1].(*object.Float)
			if !okX || !okY {
				// Try converting integers to floats
				var valX, valY float64
				if okX {
					valX = x.Value
				} else if intX, ok := args[0].(*object.Integer); ok {
					valX = float64(intX.Value)
				} else {
					return newError("arguments to `pow` must be numbers, got %s and %s", args[0].Type(), args[1].Type())
				}
				if okY {
					valY = y.Value
				} else if intY, ok := args[1].(*object.Integer); ok {
					valY = float64(intY.Value)
				} else {
					return newError("arguments to `pow` must be numbers, got %s and %s", args[0].Type(), args[1].Type())
				}
				return &object.Float{Value: math.Pow(valX, valY)}
			}
			return &object.Float{Value: math.Pow(x.Value, y.Value)}
		},
	}

	// sin
	builtins["sin"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Float:
				return &object.Float{Value: math.Sin(arg.Value)}
			case *object.Integer:
				return &object.Float{Value: math.Sin(float64(arg.Value))}
			default:
				return newError("argument to `sin` must be numeric, got %s", args[0].Type())
			}
		},
	}

	// cos
	builtins["cos"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Float:
				return &object.Float{Value: math.Cos(arg.Value)}
			case *object.Integer:
				return &object.Float{Value: math.Cos(float64(arg.Value))}
			default:
				return newError("argument to `cos` must be numeric, got %s", args[0].Type())
			}
		},
	}

	// tan
	builtins["tan"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Float:
				return &object.Float{Value: math.Tan(arg.Value)}
			case *object.Integer:
				return &object.Float{Value: math.Tan(float64(arg.Value))}
			default:
				return newError("argument to `tan` must be numeric, got %s", args[0].Type())
			}
		},
	}

	// contains
	builtins["contains"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, okS := args[0].(*object.String)
			sub, okSub := args[1].(*object.String)
			if !okS || !okSub {
				return newError("arguments to `contains` must be strings, got %s and %s", args[0].Type(), args[1].Type())
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	}

	// split
	builtins["split"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, okS := args[0].(*object.String)
			sep, okSep := args[1].(*object.String)
			if !okS || !okSep {
				return newError("arguments to `split` must be strings, got %s and %s", args[0].Type(), args[1].Type())
			}
			parts := strings.Split(s.Value, sep.Value)
			var elements []object.Object
			for _, p := range parts {
				elements = append(elements, &object.String{Value: p})
			}
			return &object.Array{Elements: elements}
		},
	}

	// join
	builtins["join"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, okArr := args[0].(*object.Array)
			sep, okSep := args[1].(*object.String)
			if !okArr || !okSep {
				return newError("arguments to `join` must be [array, string], got %s and %s", args[0].Type(), args[1].Type())
			}
			var strs []string
			for _, el := range arr.Elements {
				strs = append(strs, el.Inspect())
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	}

	// to_upper
	builtins["to_upper"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	}

	// to_lower
	builtins["to_lower"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	}

	// sleep
	builtins["sleep"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			dur, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be integer milliseconds, got %s", args[0].Type())
			}
			time.Sleep(time.Duration(dur.Value) * time.Millisecond)
			return NULL
		},
	}

	// now
	builtins["now"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	}

	// read_file
	builtins["read_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be string, got %s", args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("failed to read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	}

	// write_file
	builtins["write_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, okPath := args[0].(*object.String)
			content, okContent := args[1].(*object.String)
			if !okPath || !okContent {
				return newError("arguments to `write_file` must be [string, string], got %s and %s", args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file: %s", err.Error())
			}
			return NULL
		},
	}
}

var builtins = map[string]*object.Builtin{}
