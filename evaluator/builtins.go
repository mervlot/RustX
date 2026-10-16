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
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.Inspect())
			}
			fmt.Print(strings.Join(strArgs, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var strArgs []string
			for _, arg := range args {
				strArgs = append(strArgs, arg.Inspect())
			}
			fmt.Println(strings.Join(strArgs, " "))
			return NULL
		},
	},
	// String built-ins
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be (string, string)")
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `split` must be (string, string)")
			}
			parts := strings.Split(str.Value, sep.Value)
			var elems []object.Object
			for _, p := range parts {
				elems = append(elems, &object.String{Value: p})
			}
			return &object.Array{Elements: elems}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `join` must be (array, string)")
			}
			var strElems []string
			for _, el := range arr.Elements {
				strElems = append(strElems, el.Inspect())
			}
			return &object.String{Value: strings.Join(strElems, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
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
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	// Math built-ins
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
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
				return newError("argument to `abs` must be int or float")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			baseVal, isInt1 := args[0].(*object.Integer)
			expVal, isInt2 := args[1].(*object.Integer)
			if isInt1 && isInt2 {
				return &object.Integer{Value: int64(math.Pow(float64(baseVal.Value), float64(expVal.Value)))}
			}
			var b, e float64
			if isInt1 {
				b = float64(baseVal.Value)
			} else if f, ok := args[0].(*object.Float); ok {
				b = f.Value
			} else {
				return newError("arguments to `pow` must be numbers")
			}

			if isInt2 {
				e = float64(expVal.Value)
			} else if f, ok := args[1].(*object.Float); ok {
				e = f.Value
			} else {
				return newError("arguments to `pow` must be numbers")
			}

			return &object.Float{Value: math.Pow(b, e)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := getFloatVal(args[0])
			if math.IsNaN(val) {
				return newError("argument to `sin` must be int or float")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := getFloatVal(args[0])
			if math.IsNaN(val) {
				return newError("argument to `cos` must be int or float")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := getFloatVal(args[0])
			if math.IsNaN(val) {
				return newError("argument to `tan` must be int or float")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	// Time built-ins
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be int (milliseconds)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	// Filesystem built-ins
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be string")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("read_file error: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be (string, string)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("write_file error: %s", err.Error())
			}
			return TRUE
		},
	},
}

func getFloatVal(obj object.Object) float64 {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	default:
		return math.NaN()
	}
}
