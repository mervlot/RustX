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

	// Math built-ins
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("abs expects exactly 1 argument, got %d", len(args))
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
				return newBuiltinError("argument to abs must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("pow expects exactly 2 arguments, got %d", len(args))
			}
			var base, exponent float64

			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return newBuiltinError("first argument to pow must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			switch e := args[1].(type) {
			case *object.Integer:
				exponent = float64(e.Value)
			case *object.Float:
				exponent = e.Value
			default:
				return newBuiltinError("second argument to pow must be INTEGER or FLOAT, got %s", args[1].Type())
			}

			return &object.Float{Value: math.Pow(base, exponent)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("sin expects exactly 1 argument, got %d", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newBuiltinError("argument to sin must be INTEGER or FLOAT, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("cos expects exactly 1 argument, got %d", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newBuiltinError("argument to cos must be INTEGER or FLOAT, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("tan expects exactly 1 argument, got %d", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newBuiltinError("argument to tan must be INTEGER or FLOAT, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},

	// Strings/Arrays basic functions
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("len expects exactly 1 argument, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newBuiltinError("argument to len not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("contains expects exactly 2 arguments, got %d", len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to contains must be STRINGS, got %s and %s", args[0].Type(), args[1].Type())
			}
			return &object.Boolean{Value: strings.Contains(s.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("split expects exactly 2 arguments, got %d", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to split must be STRINGS, got %s and %s", args[0].Type(), args[1].Type())
			}
			parts := strings.Split(s.Value, sep.Value)
			var elems []object.Object
			for _, part := range parts {
				elems = append(elems, &object.String{Value: part})
			}
			return &object.Array{Elements: elems}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("join expects exactly 2 arguments, got %d", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to join must be ARRAY and STRING, got %s and %s", args[0].Type(), args[1].Type())
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
				return newBuiltinError("to_upper expects exactly 1 argument, got %d", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to to_upper must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("to_lower expects exactly 1 argument, got %d", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to to_lower must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},

	// Time functions
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("sleep expects exactly 1 argument, got %d", len(args))
			}
			var ms int64
			switch arg := args[0].(type) {
			case *object.Integer:
				ms = arg.Value
			default:
				return newBuiltinError("argument to sleep must be INTEGER (milliseconds), got %s", args[0].Type())
			}
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newBuiltinError("now expects 0 arguments, got %d", len(args))
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("read_file expects exactly 1 argument, got %d", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newBuiltinError("argument to read_file must be STRING, got %s", args[0].Type())
			}
			bytes, err := os.ReadFile(path.Value)
			if err != nil {
				return newBuiltinError("read_file failed: %s", err.Error())
			}
			return &object.String{Value: string(bytes)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("write_file expects exactly 2 arguments, got %d", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to write_file must be STRINGS (path, content), got %s and %s", args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newBuiltinError("write_file failed: %s", err.Error())
			}
			return &object.Boolean{Value: true}
		},
	},
}

func newBuiltinError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
