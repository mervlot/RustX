package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"rustx/object"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				if arg != nil {
					fmt.Print(arg.Inspect())
				}
			}
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				if arg != nil {
					out = append(out, arg.Inspect())
				}
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
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

	// String helpers
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to str_contains. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_contains must be STRING")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	"str_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to str_upper. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_upper must be STRING")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to str_lower. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_lower must be STRING")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to str_split. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_split must be STRING")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to str_join. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_join must be ARRAY and STRING")
			}
			parts := make([]string, len(arr.Elements))
			for i, elem := range arr.Elements {
				parts[i] = elem.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},

	// Math helpers
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to math_abs. got=%d, want=1", len(args))
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
				return newError("argument to math_abs must be INTEGER or FLOAT")
			}
		},
	},
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to math_sqrt. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0])
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to math_pow. got=%d, want=2", len(args))
			}
			b, err1 := getFloatVal(args[0])
			if err1 != nil {
				return err1
			}
			e, err2 := getFloatVal(args[1])
			if err2 != nil {
				return err2
			}
			return &object.Float{Value: math.Pow(b, e)}
		},
	},
	"math_floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to math_floor. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0])
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Floor(val)}
		},
	},
	"math_ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to math_ceil. got=%d, want=1", len(args))
			}
			val, errObj := getFloatVal(args[0])
			if errObj != nil {
				return errObj
			}
			return &object.Float{Value: math.Ceil(val)}
		},
	},

	// Time helpers
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to time_sleep. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to time_sleep must be INTEGER (ms)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem helpers
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to read_file. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to read_file must be STRING")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("failed to read file %s: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to write_file. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to write_file must be STRING")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file %s: %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
	"file_exists": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to file_exists. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to file_exists must be STRING")
			}
			_, err := os.Stat(path.Value)
			return nativeBoolToBooleanObject(!os.IsNotExist(err))
		},
	},
}

func getFloatVal(obj object.Object) (float64, object.Object) {
	switch o := obj.(type) {
	case *object.Float:
		return o.Value, nil
	case *object.Integer:
		return float64(o.Value), nil
	default:
		return 0, newError("expected INTEGER or FLOAT, got %s", obj.Type())
	}
}
