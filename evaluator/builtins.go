package evaluator

import (
	"fmt"
	"math"
	"os"

	"github.com/username/rustx/object"
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
			return &object.Null{}
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
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newBuiltinError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `contains` must be strings")
			}
			for i := 0; i <= len(str.Value)-len(substr.Value); i++ {
				if str.Value[i:i+len(substr.Value)] == substr.Value {
					return &object.Boolean{Value: true}
				}
			}
			return &object.Boolean{Value: false}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `split` must be strings")
			}
			var result []object.Object
			if sep.Value == "" {
				for _, ch := range str.Value {
					result = append(result, &object.String{Value: string(ch)})
				}
			} else {
				start := 0
				for i := 0; i <= len(str.Value)-len(sep.Value); i++ {
					if str.Value[i:i+len(sep.Value)] == sep.Value {
						result = append(result, &object.String{Value: str.Value[start:i]})
						start = i + len(sep.Value)
						i += len(sep.Value) - 1
					}
				}
				result = append(result, &object.String{Value: str.Value[start:]})
			}
			return &object.Array{Elements: result}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("first argument to `join` must be Array, second must be String")
			}
			var str string
			for i, el := range arr.Elements {
				if i > 0 {
					str += sep.Value
				}
				str += el.Inspect()
			}
			return &object.String{Value: str}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `to_upper` must be String")
			}
			var res []byte
			for i := 0; i < len(str.Value); i++ {
				ch := str.Value[i]
				if 'a' <= ch && ch <= 'z' {
					res = append(res, ch-('a'-'A'))
				} else {
					res = append(res, ch)
				}
			}
			return &object.String{Value: string(res)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `to_lower` must be String")
			}
			var res []byte
			for i := 0; i < len(str.Value); i++ {
				ch := str.Value[i]
				if 'A' <= ch && ch <= 'Z' {
					res = append(res, ch+('a'-'A'))
				} else {
					res = append(res, ch)
				}
			}
			return &object.String{Value: string(res)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
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
				return newBuiltinError("argument to `abs` must be number")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			b, ok1 := getFloatValue(args[0])
			exp, ok2 := getFloatValue(args[1])
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `pow` must be numbers")
			}
			res := math.Pow(b, exp)
			if _, isInt := args[0].(*object.Integer); isInt {
				if _, isInt2 := args[1].(*object.Integer); isInt2 {
					return &object.Integer{Value: int64(res)}
				}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newBuiltinError("argument to `sin` must be number")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newBuiltinError("argument to `cos` must be number")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newBuiltinError("argument to `tan` must be number")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			filename, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `read_file` must be String")
			}
			content, err := os.ReadFile(filename.Value)
			if err != nil {
				return newBuiltinError("failed to read file %s: %s", filename.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			filename, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `write_file` must be Strings")
			}
			err := os.WriteFile(filename.Value, []byte(content.Value), 0644)
			if err != nil {
				return newBuiltinError("failed to write file %s: %s", filename.Value, err.Error())
			}
			return &object.Boolean{Value: true}
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

func newBuiltinError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
