package evaluator

import (
	"fmt"
	"math"

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
				return newBuiltinError("len() expects exactly 1 argument, got %d", len(args))
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
				return newBuiltinError("contains() expects 2 arguments")
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("contains() arguments must be strings")
			}
			return &object.Boolean{Value: containsSubstring(s.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("split() expects 2 arguments")
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("split() arguments must be strings")
			}
			parts := splitString(s.Value, sep.Value)
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
				return newBuiltinError("join() expects 2 arguments")
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("join() arguments must be (array, string)")
			}
			strs := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strs[i] = el.Inspect()
			}
			return &object.String{Value: joinStrings(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("to_upper() expects 1 argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("to_upper() argument must be string")
			}
			return &object.String{Value: toUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("to_lower() expects 1 argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("to_lower() argument must be string")
			}
			return &object.String{Value: toLower(s.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("abs() expects 1 argument")
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
				return newBuiltinError("abs() argument must be integer or float")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("pow() expects 2 arguments")
			}
			base := getFloatValue(args[0])
			exp := getFloatValue(args[1])
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("sin() expects 1 argument")
			}
			val := getFloatValue(args[0])
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("cos() expects 1 argument")
			}
			val := getFloatValue(args[0])
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("tan() expects 1 argument")
			}
			val := getFloatValue(args[0])
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: getCurrentUnixTimestamp()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("sleep() expects 1 argument (milliseconds)")
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newBuiltinError("sleep() argument must be integer milliseconds")
			}
			sleepMilliseconds(ms.Value)
			return &object.Null{}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("read_file() expects 1 argument")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("read_file() argument must be string path")
			}
			content, err := readFileContent(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error())}
			}
			return &object.String{Value: content}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("write_file() expects 2 arguments (path, content)")
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("write_file() arguments must be strings")
			}
			err := writeFileContent(path.Value, content.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func newBuiltinError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func getFloatValue(obj object.Object) float64 {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	default:
		return 0
	}
}
