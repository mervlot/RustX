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
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError(pos, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(pos, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(pos, "arguments to `contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(pos, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(pos, "arguments to `split` must be strings")
			}
			parts := strings.Split(str.Value, sep.Value)
			// Return space joined for simple script manipulation
			return &object.String{Value: strings.Join(parts, ",")}
		},
	},
	"join": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(pos, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(pos, "arguments to `join` must be strings")
			}
			parts := strings.Split(str.Value, ",")
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(pos, "argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(pos, "argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
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
				return newError(pos, "argument to `abs` must be int or float")
			}
		},
	},
	"pow": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(pos, "wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, exponent float64
			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return newError(pos, "first argument to `pow` must be number")
			}
			switch e := args[1].(type) {
			case *object.Integer:
				exponent = float64(e.Value)
			case *object.Float:
				exponent = e.Value
			default:
				return newError(pos, "second argument to `pow` must be number")
			}
			return &object.Float{Value: math.Pow(base, exponent)}
		},
	},
	"sin": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch a := args[0].(type) {
			case *object.Integer:
				val = float64(a.Value)
			case *object.Float:
				val = a.Value
			default:
				return newError(pos, "argument to `sin` must be number")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch a := args[0].(type) {
			case *object.Integer:
				val = float64(a.Value)
			case *object.Float:
				val = a.Value
			default:
				return newError(pos, "argument to `cos` must be number")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch a := args[0].(type) {
			case *object.Integer:
				val = float64(a.Value)
			case *object.Float:
				val = a.Value
			default:
				return newError(pos, "argument to `tan` must be number")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(pos, "argument to `sleep` must be integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixMilli()}
		},
	},
	"read_file": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(pos, "wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(pos, "argument to `read_file` must be string")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(pos, "failed to read file '%s': %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(pos, "wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(pos, "arguments to `write_file` must be strings")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(pos, "failed to write file '%s': %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}
