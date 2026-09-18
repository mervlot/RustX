package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

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
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(0, 0, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `contains`. got=%d, want=2", len(args))
			}
			switch s := args[0].(type) {
			case *object.String:
				sub, ok := args[1].(*object.String)
				if !ok {
					return newError(0, 0, "second argument to `contains` for string must be string, got %s", args[1].Type())
				}
				return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
			case *object.Array:
				for _, elem := range s.Elements {
					if evalObjectEquality(elem, args[1]) {
						return TRUE
					}
				}
				return FALSE
			default:
				return newError(0, 0, "first argument to `contains` must be string or array, got %s", args[0].Type())
			}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `split`. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `split` must be strings")
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
				return newError(0, 0, "wrong number of arguments to `join`. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "first argument to `join` must be array and second string")
			}
			strParts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strParts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(strParts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `to_upper`. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_upper` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `to_lower`. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_lower` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `abs`. got=%d, want=1", len(args))
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
				return newError(0, 0, "argument to `abs` must be int or float, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `pow`. got=%d, want=2", len(args))
			}
			var base, exp float64
			if b, ok := args[0].(*object.Integer); ok {
				base = float64(b.Value)
			} else if b, ok := args[0].(*object.Float); ok {
				base = b.Value
			} else {
				return newError(0, 0, "arguments to `pow` must be numbers")
			}

			if e, ok := args[1].(*object.Integer); ok {
				exp = float64(e.Value)
			} else if e, ok := args[1].(*object.Float); ok {
				exp = e.Value
			} else {
				return newError(0, 0, "arguments to `pow` must be numbers")
			}

			res := math.Pow(base, exp)
			if _, ok1 := args[0].(*object.Integer); ok1 {
				if _, ok2 := args[1].(*object.Integer); ok2 {
					return &object.Integer{Value: int64(res)}
				}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			val, err := getFloatArg("sin", args)
			if err != nil {
				return err
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			val, err := getFloatArg("cos", args)
			if err != nil {
				return err
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			val, err := getFloatArg("tan", args)
			if err != nil {
				return err
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `sleep`. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(0, 0, "argument to `sleep` must be integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `read_file`. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `read_file` must be string path")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(0, 0, "failed to read file '%s': %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `write_file`. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `write_file` must be strings (path, content)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(0, 0, "failed to write file '%s': %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func getFloatArg(fnName string, args []object.Object) (float64, object.Object) {
	if len(args) != 1 {
		return 0, newError(0, 0, "wrong number of arguments to `%s`. got=%d, want=1", fnName, len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return float64(arg.Value), nil
	case *object.Float:
		return arg.Value, nil
	default:
		return 0, newError(0, 0, "argument to `%s` must be number, got %s", fnName, args[0].Type())
	}
}
