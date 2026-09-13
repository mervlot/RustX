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
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Print(strings.Join(parts, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(strings.Join(parts, " "))
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
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_upper` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `to_lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `contains` must be STRING")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `split` must be STRING")
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
				return newError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `join` must be ARRAY and STRING")
			}
			var parts []string
			for _, elem := range arr.Elements {
				parts = append(parts, elem.Inspect())
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
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
				return newError(0, 0, "argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
			}
			baseVal, isNum1 := getFloatValue(args[0])
			expVal, isNum2 := getFloatValue(args[1])
			if !isNum1 || !isNum2 {
				return newError(0, 0, "arguments to `pow` must be numeric")
			}
			res := math.Pow(baseVal, expVal)
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(0, 0, "argument to `sin` must be numeric")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(0, 0, "argument to `cos` must be numeric")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(0, 0, "argument to `tan` must be numeric")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(0, 0, "argument to `sleep` must be INTEGER (ms)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to `read_file` must be STRING")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(0, 0, "could not read file %s: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to `write_file` must be STRING")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(0, 0, "could not write file %s: %s", path.Value, err.Error())
			}
			return nativeBoolToBooleanObject(true)
		},
	},
}

func getFloatValue(obj object.Object) (float64, bool) {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value), true
	case *object.Float:
		return v.Value, true
	default:
		return 0, false
	}
}
