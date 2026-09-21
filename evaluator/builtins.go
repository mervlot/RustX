package evaluator

import (
	"fmt"
	"math"

	"github.com/username/rustx/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var builtins = map[string]*object.Builtin{
	// Printing
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
			return NULL
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
			return NULL
		},
	},
}

func init() {
	// Register standard library functions into builtins map
	registerStdLib()
}

func registerStdLib() {
	// Strings stdlib
	builtins["len"] = &object.Builtin{
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
	}

	builtins["contains"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be STRING, STRING")
			}
			if containsString(s.Value, substr.Value) {
				return TRUE
			}
			return FALSE
		},
	}

	builtins["split"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `split` must be STRING, STRING")
			}
			parts := splitString(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	}

	builtins["join"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `join` must be ARRAY, STRING")
			}
			strs := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strs[i] = el.Inspect()
			}
			return &object.String{Value: joinStrings(strs, sep.Value)}
		},
	}

	builtins["to_upper"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be STRING")
			}
			return &object.String{Value: toUpperString(s.Value)}
		},
	}

	builtins["to_lower"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be STRING")
			}
			return &object.String{Value: toLowerString(s.Value)}
		},
	}

	// Math stdlib
	builtins["abs"] = &object.Builtin{
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
				return newError("argument to `abs` must be INTEGER or FLOAT")
			}
		},
	}

	builtins["pow"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			base := getFloatVal(args[0])
			exponent := getFloatVal(args[1])
			if base == nil || exponent == nil {
				return newError("arguments to `pow` must be numbers")
			}
			return &object.Float{Value: math.Pow(*base, *exponent)}
		},
	}

	builtins["sin"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := getFloatVal(args[0])
			if v == nil {
				return newError("argument to `sin` must be number")
			}
			return &object.Float{Value: math.Sin(*v)}
		},
	}

	builtins["cos"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := getFloatVal(args[0])
			if v == nil {
				return newError("argument to `cos` must be number")
			}
			return &object.Float{Value: math.Cos(*v)}
		},
	}

	builtins["tan"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := getFloatVal(args[0])
			if v == nil {
				return newError("argument to `tan` must be number")
			}
			return &object.Float{Value: math.Tan(*v)}
		},
	}

	// Time stdlib
	builtins["now"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: getUnixTime()}
		},
	}

	builtins["sleep"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be INTEGER (milliseconds)")
			}
			sleepMs(ms.Value)
			return NULL
		},
	}

	// Filesystem stdlib
	builtins["read_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be STRING")
			}
			content, err := readFileContent(path.Value)
			if err != nil {
				return newError("failed to read file %s: %s", path.Value, err.Error())
			}
			return &object.String{Value: content}
		},
	}

	builtins["write_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be STRING, STRING")
			}
			err := writeFileContent(path.Value, content.Value)
			if err != nil {
				return newError("failed to write file %s: %s", path.Value, err.Error())
			}
			return TRUE
		},
	}
}

func getFloatVal(obj object.Object) *float64 {
	switch o := obj.(type) {
	case *object.Integer:
		f := float64(o.Value)
		return &f
	case *object.Float:
		return &o.Value
	default:
		return nil
	}
}
