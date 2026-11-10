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
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newBuiltinError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newBuiltinError("wrong number of arguments to append. got=%d, want>=2", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newBuiltinError("first argument to `append` must be ARRAY, got %s", args[0].Type())
			}
			newElements := make([]object.Object, len(arr.Elements), len(arr.Elements)+len(args)-1)
			copy(newElements, arr.Elements)
			newElements = append(newElements, args[1:]...)
			return &object.Array{Elements: newElements}
		},
	},

	// Strings module / built-ins
	"str_len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `str_len` must be STRING, got %s", args[0].Type())
			}
			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `to_upper` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `to_lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `contains` must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `split` must be STRING")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},

	// Math built-ins
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments to `abs`. got=%d, want=1", len(args))
			}
			switch val := args[0].(type) {
			case *object.Integer:
				if val.Value < 0 {
					return &object.Integer{Value: -val.Value}
				}
				return val
			case *object.Float:
				return &object.Float{Value: math.Abs(val.Value)}
			default:
				return newBuiltinError("argument to `abs` must be INT or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments to `sqrt`. got=%d, want=1", len(args))
			}
			var num float64
			switch val := args[0].(type) {
			case *object.Integer:
				num = float64(val.Value)
			case *object.Float:
				num = val.Value
			default:
				return newBuiltinError("argument to `sqrt` must be INT or FLOAT, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sqrt(num)}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments to `pow`. got=%d, want=2", len(args))
			}
			var x, y float64
			switch v := args[0].(type) {
			case *object.Integer:
				x = float64(v.Value)
			case *object.Float:
				x = v.Value
			default:
				return newBuiltinError("first argument to `pow` must be numeric")
			}
			switch v := args[1].(type) {
			case *object.Integer:
				y = float64(v.Value)
			case *object.Float:
				y = v.Value
			default:
				return newBuiltinError("second argument to `pow` must be numeric")
			}
			return &object.Float{Value: math.Pow(x, y)}
		},
	},
	"min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments to `min`. got=%d, want=2", len(args))
			}
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				v1 := args[0].(*object.Integer).Value
				v2 := args[1].(*object.Integer).Value
				if v1 < v2 {
					return args[0]
				}
				return args[1]
			}
			f1, err1 := toFloat(args[0])
			f2, err2 := toFloat(args[1])
			if err1 != nil || err2 != nil {
				return newBuiltinError("arguments to `min` must be numeric")
			}
			return &object.Float{Value: math.Min(f1, f2)}
		},
	},
	"max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments to `max`. got=%d, want=2", len(args))
			}
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				v1 := args[0].(*object.Integer).Value
				v2 := args[1].(*object.Integer).Value
				if v1 > v2 {
					return args[0]
				}
				return args[1]
			}
			f1, err1 := toFloat(args[0])
			f2, err2 := toFloat(args[1])
			if err1 != nil || err2 != nil {
				return newBuiltinError("arguments to `max` must be numeric")
			}
			return &object.Float{Value: math.Max(f1, f2)}
		},
	},

	// Time built-ins
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments to `time_sleep`. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newBuiltinError("argument to `time_sleep` must be INT (milliseconds)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments to `read_file`. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to `read_file` must be STRING")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newBuiltinError("failed to read file %q: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments to `write_file`. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `write_file` must be (STRING, STRING)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newBuiltinError("failed to write file %q: %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func newBuiltinError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func toFloat(obj object.Object) (float64, error) {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value), nil
	case *object.Float:
		return v.Value, nil
	default:
		return 0, fmt.Errorf("not a number")
	}
}
