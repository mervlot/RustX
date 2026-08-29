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
	// Printing
	"println": {
		Fn: func(args ...object.Object) object.Object {
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(out, " "))
			return &object.Null{}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Print(strings.Join(out, " "))
			return &object.Null{}
		},
	},

	// Strings / Collection
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
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to contains. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to contains must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to split. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to split must be strings")
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
				return newError("wrong number of arguments to join. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to join must be array and string")
			}
			strs := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strs[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to to_upper. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to to_upper must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to to_lower. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to to_lower must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},

	// Math builtins
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to abs. got=%d, want=1", len(args))
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
				return newError("argument to abs must be int or float")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to pow. got=%d, want=2", len(args))
			}
			base := objectToFloat(args[0])
			exp := objectToFloat(args[1])
			if base == nil || exp == nil {
				return newError("arguments to pow must be numbers")
			}
			return &object.Float{Value: math.Pow(*base, *exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to sin. got=%d, want=1", len(args))
			}
			v := objectToFloat(args[0])
			if v == nil {
				return newError("argument to sin must be number")
			}
			return &object.Float{Value: math.Sin(*v)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to cos. got=%d, want=1", len(args))
			}
			v := objectToFloat(args[0])
			if v == nil {
				return newError("argument to cos must be number")
			}
			return &object.Float{Value: math.Cos(*v)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to tan. got=%d, want=1", len(args))
			}
			v := objectToFloat(args[0])
			if v == nil {
				return newError("argument to tan must be number")
			}
			return &object.Float{Value: math.Tan(*v)}
		},
	},

	// Time builtins
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to sleep. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to sleep must be int (milliseconds)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},

	// Filesystem builtins
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments to read_file. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to read_file must be string")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("could not read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
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
				return newError("arguments to write_file must be strings (path, content)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("could not write file: %s", err.Error())
			}
			return &object.Null{}
		},
	},
}

func objectToFloat(obj object.Object) *float64 {
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
