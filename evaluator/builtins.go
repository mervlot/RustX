package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(tok, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to `contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to `split` must be strings")
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
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "first argument to `join` must be array, second must be string")
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
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
				return newError(tok, "argument to `abs` must be integer or float")
			}
		},
	},
	"pow": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments. got=%d, want=2", len(args))
			}
			base, bOk := getFloatValue(args[0])
			exp, eOk := getFloatValue(args[1])
			if !bOk || !eOk {
				return newError(tok, "arguments to `pow` must be numbers")
			}
			res := math.Pow(base, exp)
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to `sin` must be number")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to `cos` must be number")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to `tan` must be number")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(tok, "argument to `sleep` must be integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError(tok, "wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Integer{Value: time.Now().UnixMilli()}
		},
	},
	"read_file": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to `read_file` must be string path")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(tok, "could not read file: %s", err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to `write_file` must be strings")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(tok, "could not write file: %s", err.Error())
			}
			return TRUE
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
