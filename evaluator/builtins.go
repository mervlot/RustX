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
	// Printing
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

	// Strings / Collection Builtins
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
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `split`. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `split` must be strings")
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
				return newError(0, 0, "wrong number of arguments to `join`. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `join` must be (array, string)")
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
				return newError(0, 0, "wrong number of arguments to `to_upper`. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `to_lower`. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},

	// Math Builtins
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
				return newError(0, 0, "argument to `abs` must be number, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `pow`. got=%d, want=2", len(args))
			}
			var base, exp float64
			if bInt, ok := args[0].(*object.Integer); ok {
				base = float64(bInt.Value)
			} else if bFlt, ok := args[0].(*object.Float); ok {
				base = bFlt.Value
			} else {
				return newError(0, 0, "first argument to `pow` must be number")
			}

			if eInt, ok := args[1].(*object.Integer); ok {
				exp = float64(eInt.Value)
			} else if eFlt, ok := args[1].(*object.Float); ok {
				exp = eFlt.Value
			} else {
				return newError(0, 0, "second argument to `pow` must be number")
			}

			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `sin`. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0], "sin")
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `cos`. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0], "cos")
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `tan`. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0], "tan")
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},

	// Time Builtins
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `sleep`. got=%d, want=1", len(args))
			}
			i, ok := args[0].(*object.Integer)
			if !ok {
				return newError(0, 0, "argument to `sleep` must be integer milliseconds")
			}
			time.Sleep(time.Duration(i.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem Basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments to `read_file`. got=%d, want=1", len(args))
			}
			filename, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `read_file` must be string filename")
			}
			content, err := os.ReadFile(filename.Value)
			if err != nil {
				return newError(0, 0, "failed to read file '%s': %s", filename.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments to `write_file`. got=%d, want=2", len(args))
			}
			filename, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `write_file` must be (filename, content) strings")
			}
			err := os.WriteFile(filename.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(0, 0, "failed to write file '%s': %s", filename.Value, err.Error())
			}
			return TRUE
		},
	},
}

func getFloatVal(obj object.Object, fnName string) (float64, object.Object) {
	if i, ok := obj.(*object.Integer); ok {
		return float64(i.Value), nil
	}
	if f, ok := obj.(*object.Float); ok {
		return f.Value, nil
	}
	return 0, newError(0, 0, "argument to `%s` must be number, got %s", fnName, obj.Type())
}
