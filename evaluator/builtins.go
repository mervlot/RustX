package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"rustx/object"
)

var builtins = map[string]*object.Builtin{
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
	// Strings std
	"str_to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"str_to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments must be STRING")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, sub.Value))
		},
	},
	// Math std
	"sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `sqrt` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			getVal := func(o object.Object) (float64, bool) {
				if i, ok := o.(*object.Integer); ok {
					return float64(i.Value), true
				}
				if f, ok := o.(*object.Float); ok {
					return f.Value, true
				}
				return 0, false
			}
			x, ok1 := getVal(args[0])
			y, ok2 := getVal(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `pow` must be numbers")
			}
			return &object.Float{Value: math.Pow(x, y)}
		},
	},
	// Time std
	"now_ms": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"sleep_ms": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep_ms` must be INTEGER, got %s", args[0].Type())
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
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
				return newError("argument to `read_file` must be STRING, got %s", args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("failed to read file %s: %s", path.Value, err.Error())
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
				return newError("arguments to `write_file` must be STRING")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file %s: %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}
