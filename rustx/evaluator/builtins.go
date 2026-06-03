package evaluator

import (
	"fmt"
	"github.com/username/rustx/object"
	"math"
	"os"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	// Math
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
				return newError("argument to `abs` not supported, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			var x, y float64
			if arg, ok := args[0].(*object.Float); ok {
				x = arg.Value
			} else if arg, ok := args[0].(*object.Integer); ok {
				x = float64(arg.Value)
			} else {
				return newError("argument to `pow` not supported, got %s", args[0].Type())
			}
			if arg, ok := args[1].(*object.Float); ok {
				y = arg.Value
			} else if arg, ok := args[1].(*object.Integer); ok {
				y = float64(arg.Value)
			} else {
				return newError("argument to `pow` not supported, got %s", args[1].Type())
			}
			return &object.Float{Value: math.Pow(x, y)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			if arg, ok := args[0].(*object.Float); ok {
				val = arg.Value
			} else if arg, ok := args[0].(*object.Integer); ok {
				val = float64(arg.Value)
			} else {
				return newError("argument to `sin` not supported, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			if arg, ok := args[0].(*object.Float); ok {
				val = arg.Value
			} else if arg, ok := args[0].(*object.Integer); ok {
				val = float64(arg.Value)
			} else {
				return newError("argument to `cos` not supported, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	// Strings
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `split` must be strings")
			}
			parts := strings.Split(s.Value, sep.Value)
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
				return newError("arguments to `join` must be (Array, String)")
			}
			parts := make([]string, len(arr.Elements))
			for i, e := range arr.Elements {
				parts[i] = e.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	// Time
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be integer (ms)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	// Filesystem
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
				return newError("could not read file: %s", err)
			}
			return &object.String{Value: string(content)}
		},
	},
}
