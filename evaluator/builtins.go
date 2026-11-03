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

	// Strings
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
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be STRING, got %s", args[0].Type())
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
				return newError("argument to `to_lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, sub.Value))
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
				return newError("arguments to `join` must be ARRAY and STRING")
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `trim` must be STRING")
			}
			return &object.String{Value: strings.TrimSpace(str.Value)}
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
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Integer{Value: val}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to `abs` must be INTEGER or FLOAT")
			}
		},
	},
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
				return newError("argument to `sqrt` must be number")
			}
			if val < 0 {
				return newError("cannot take square root of negative number: %g", val)
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if a, ok1 := args[0].(*object.Integer); ok1 {
				if b, ok2 := args[1].(*object.Integer); ok2 {
					if a.Value < b.Value {
						return a
					}
					return b
				}
			}
			f1, ok1 := toFloat(args[0])
			f2, ok2 := toFloat(args[1])
			if ok1 && ok2 {
				return &object.Float{Value: math.Min(f1, f2)}
			}
			return newError("arguments to `min` must be numbers")
		},
	},
	"max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if a, ok1 := args[0].(*object.Integer); ok1 {
				if b, ok2 := args[1].(*object.Integer); ok2 {
					if a.Value > b.Value {
						return a
					}
					return b
				}
			}
			f1, ok1 := toFloat(args[0])
			f2, ok2 := toFloat(args[1])
			if ok1 && ok2 {
				return &object.Float{Value: math.Max(f1, f2)}
			}
			return newError("arguments to `max` must be numbers")
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			f1, ok1 := toFloat(args[0])
			f2, ok2 := toFloat(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `pow` must be numbers")
			}
			return &object.Float{Value: math.Pow(f1, f2)}
		},
	},
	"floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `floor` must be number")
			}
			return &object.Float{Value: math.Floor(f)}
		},
	},
	"ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `ceil` must be number")
			}
			return &object.Float{Value: math.Ceil(f)}
		},
	},
	"round": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `round` must be number")
			}
			return &object.Float{Value: math.Round(f)}
		},
	},

	// Time
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be INTEGER milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// File System Basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be STRING")
			}
			content, err := os.ReadFile(filepath.Value)
			if err != nil {
				return newError("failed to read file '%s': %s", filepath.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			filepath, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be strings (filepath, content)")
			}
			err := os.WriteFile(filepath.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file '%s': %s", filepath.Value, err.Error())
			}
			return TRUE
		},
	},

	// Array helpers
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			newElements := append(arr.Elements, args[1])
			return &object.Array{Elements: newElements}
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `pop` must be ARRAY, got %s", args[0].Type())
			}
			length := len(arr.Elements)
			if length == 0 {
				return NULL
			}
			last := arr.Elements[length-1]
			arr.Elements = arr.Elements[:length-1]
			return last
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			length := len(arr.Elements)
			if length > 1 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return &object.Array{Elements: []object.Object{}}
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
