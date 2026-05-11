package evaluator

import (
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

func init() {
	// Add more built-ins
	builtins["len"] = &object.Builtin{
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
	}

	// Strings
	builtins["str_contains"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `str_contains` must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	}

	// Math
	builtins["math_sqrt"] = &object.Builtin{
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
				return newError("argument to `math_sqrt` must be number")
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	}

	// Time
	builtins["time_now"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.String{Value: time.Now().Format(time.RFC3339)}
		},
	}

	// Filesystem
	builtins["fs_read"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `fs_read` must be string")
			}
			content, err := ioutil.ReadFile(path.Value)
			if err != nil {
				return newError("error reading file: %s", err)
			}
			return &object.String{Value: string(content)}
		},
	}
}
