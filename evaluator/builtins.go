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
			strArgs := []string{}
			for _, arg := range args {
				strArgs = append(strArgs, arg.Inspect())
			}
			fmt.Print(strings.Join(strArgs, " "))
			return &object.Null{}
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			strArgs := []string{}
			for _, arg := range args {
				strArgs = append(strArgs, arg.Inspect())
			}
			fmt.Println(strings.Join(strArgs, " "))
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newErrorStatic("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `contains` must be strings, got %s and %s", args[0].Type(), args[1].Type())
			}
			return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `split` must be strings, got %s and %s", args[0].Type(), args[1].Type())
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
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `join` must be (array, string), got %s and %s", args[0].Type(), args[1].Type())
			}
			parts := []string{}
			for _, el := range arr.Elements {
				parts = append(parts, el.Inspect())
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `to_upper` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `to_lower` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
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
				return newErrorStatic("argument to `abs` must be int or float, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, exp float64

			switch arg := args[0].(type) {
			case *object.Integer:
				base = float64(arg.Value)
			case *object.Float:
				base = arg.Value
			default:
				return newErrorStatic("first argument to `pow` must be number, got %s", args[0].Type())
			}

			switch arg := args[1].(type) {
			case *object.Integer:
				exp = float64(arg.Value)
			case *object.Float:
				exp = arg.Value
			default:
				return newErrorStatic("second argument to `pow` must be number, got %s", args[1].Type())
			}

			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newErrorStatic("argument to `sin` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newErrorStatic("argument to `cos` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newErrorStatic("argument to `tan` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			dur, ok := args[0].(*object.Integer)
			if !ok {
				return newErrorStatic("argument to `sleep` must be integer (milliseconds), got %s", args[0].Type())
			}
			time.Sleep(time.Duration(dur.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newErrorStatic("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Float{Value: float64(time.Now().UnixNano()) / 1e9}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `read_file` must be string, got %s", args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newErrorStatic("failed to read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `write_file` must be (string, string), got %s and %s", args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Boolean{Value: false}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func newErrorStatic(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    0,
		Column:  0,
	}
}
