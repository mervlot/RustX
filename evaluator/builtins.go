package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/username/rustx/object"
)

var builtins = map[string]*object.Builtin{
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
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'len': got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to 'len' not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to 'contains': got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to 'contains' must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to 'split': got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to 'split' must be strings")
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
				return newError("wrong number of arguments to 'join': got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("first argument to 'join' must be array, second must be string")
			}
			strs := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strs[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'to_upper': got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to 'to_upper' must be string")
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'to_lower': got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to 'to_lower' must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'abs': got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				v := arg.Value
				if v < 0 {
					v = -v
				}
				return &object.Integer{Value: v}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to 'abs' must be integer or float")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to 'pow': got=%d, want=2", len(args))
			}
			base, ok1 := getFloatValue(args[0])
			exponent, ok2 := getFloatValue(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to 'pow' must be numbers")
			}
			return &object.Float{Value: math.Pow(base, exponent)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'sin': got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to 'sin' must be a number")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'cos': got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to 'cos' must be a number")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'tan': got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError("argument to 'tan' must be a number")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to 'read_file': got=%d, want=1", len(args))
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to 'read_file' must be string")
			}
			content, err := os.ReadFile(filepath.Value)
			if err != nil {
				return newError("read_file failed: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to 'write_file': got=%d, want=2", len(args))
			}
			filepath, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to 'write_file' must be strings")
			}
			err := os.WriteFile(filepath.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("write_file failed: %s", err.Error())
			}
			return nativeBoolToBooleanObject(true)
		},
	},
}

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
