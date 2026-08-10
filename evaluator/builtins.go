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
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Print(strings.Join(parts, " "))
			return &object.Null{}
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(strings.Join(parts, " "))
			return &object.Null{}
		},
	},

	// Strings
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newErrorStatic("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `contains` must be strings")
			}
			return &object.Boolean{Value: strings.Contains(s.Value, sub.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `split` must be strings")
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
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `join` must be (array, string)")
			}
			var parts []string
			for _, el := range arr.Elements {
				parts = append(parts, el.Inspect())
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `to_upper` must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `to_lower` must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},

	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Integer{Value: val}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newErrorStatic("argument to `abs` must be numerical")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			base, ok1 := args[0].(*object.Float)
			exp, ok2 := args[1].(*object.Float)
			if ok1 && ok2 {
				return &object.Float{Value: math.Pow(base.Value, exp.Value)}
			}
			// fallback to convert int if necessary
			bVal := 0.0
			if ok1 {
				bVal = base.Value
			} else if bi, ok := args[0].(*object.Integer); ok {
				bVal = float64(bi.Value)
			} else {
				return newErrorStatic("arguments to `pow` must be numerical")
			}

			eVal := 0.0
			if ok2 {
				eVal = exp.Value
			} else if ei, ok := args[1].(*object.Integer); ok {
				eVal = float64(ei.Value)
			} else {
				return newErrorStatic("arguments to `pow` must be numerical")
			}

			return &object.Float{Value: math.Pow(bVal, eVal)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := 0.0
			switch arg := args[0].(type) {
			case *object.Float:
				val = arg.Value
			case *object.Integer:
				val = float64(arg.Value)
			default:
				return newErrorStatic("argument to `sin` must be numerical")
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := 0.0
			switch arg := args[0].(type) {
			case *object.Float:
				val = arg.Value
			case *object.Integer:
				val = float64(arg.Value)
			default:
				return newErrorStatic("argument to `cos` must be numerical")
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			val := 0.0
			switch arg := args[0].(type) {
			case *object.Float:
				val = arg.Value
			case *object.Integer:
				val = float64(arg.Value)
			default:
				return newErrorStatic("argument to `tan` must be numerical")
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},

	// Time
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newErrorStatic("argument to `sleep` must be integer (milliseconds)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newErrorStatic("wrong number of arguments. got=%d, want=0", len(args))
			}
			// return current epoch millisecond
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorStatic("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newErrorStatic("argument to `read_file` must be string (filepath)")
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newErrorStatic("failed to read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorStatic("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorStatic("arguments to `write_file` must be (filepath, content) as strings")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newErrorStatic("failed to write file: %s", err.Error())
			}
			return &object.Null{}
		},
	},
}

func newErrorStatic(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}
