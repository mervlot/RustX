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
			return &object.Null{}
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
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for len. got=%d, want=1", len(args))
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
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for split. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `split` must be (string, string)")
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
				return newError("wrong number of arguments for join. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `join` must be (array, string)")
			}
			strParts := make([]string, len(arr.Elements))
			for i, elem := range arr.Elements {
				strParts[i] = elem.Inspect()
			}
			return &object.String{Value: strings.Join(strParts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for to_upper. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for to_lower. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for contains. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be (string, string)")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for abs. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				v := arg.Value
				if v < 0 {
					v = -v
				}
				return &object.Integer{Value: v}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError("argument to `abs` must be integer or float")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for pow. got=%d, want=2", len(args))
			}
			base, ok1 := toFloat(args[0])
			exp, ok2 := toFloat(args[1])
			if !ok1 || !ok2 {
				return newError("arguments to `pow` must be numbers")
			}
			res := math.Pow(base, exp)
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for sin. got=%d, want=1", len(args))
			}
			v, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `sin` must be number")
			}
			return &object.Float{Value: math.Sin(v)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for cos. got=%d, want=1", len(args))
			}
			v, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `cos` must be number")
			}
			return &object.Float{Value: math.Cos(v)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for tan. got=%d, want=1", len(args))
			}
			v, ok := toFloat(args[0])
			if !ok {
				return newError("argument to `tan` must be number")
			}
			return &object.Float{Value: math.Tan(v)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for sleep. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments for now. got=%d, want=0", len(args))
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for read_file. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be string path")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("could not read file %s: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for write_file. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be (string, string)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("could not write file %s: %s", path.Value, err.Error())
			}
			return nativeBoolToBooleanObject(true)
		},
	},
}

func toFloat(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
