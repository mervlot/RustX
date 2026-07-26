package evaluator

import (
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", line, col, args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			print(strings.Join(parts, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			println(strings.Join(parts, " "))
			return NULL
		},
	},
	"contains": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", line, col, len(args))
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `contains` must be strings, got %s and %s", line, col, args[0].Type(), args[1].Type())
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", line, col, len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `split` must be strings, got %s and %s", line, col, args[0].Type(), args[1].Type())
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
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", line, col, len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `join` must be (array, string), got %s and %s", line, col, args[0].Type(), args[1].Type())
			}
			var parts []string
			for _, el := range arr.Elements {
				parts = append(parts, el.Inspect())
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be string, got %s", line, col, args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string, got %s", line, col, args[0].Type())
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"abs": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
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
				return newError("argument to `abs` must be number, got %s", line, col, args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", line, col, len(args))
			}
			var base, exponent float64
			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return newError("arguments to `pow` must be numbers, got %s and %s", line, col, args[0].Type(), args[1].Type())
			}
			switch e := args[1].(type) {
			case *object.Integer:
				exponent = float64(e.Value)
			case *object.Float:
				exponent = e.Value
			default:
				return newError("arguments to `pow` must be numbers, got %s and %s", line, col, args[0].Type(), args[1].Type())
			}
			return &object.Float{Value: math.Pow(base, exponent)}
		},
	},
	"sin": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `sin` must be number, got %s", line, col, args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `cos` must be number, got %s", line, col, args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `tan` must be number, got %s", line, col, args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			sec, ok := args[0].(*object.Integer)
			if !ok {
				return newError("argument to `sleep` must be integer (milliseconds), got %s", line, col, args[0].Type())
			}
			time.Sleep(time.Duration(sec.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0", line, col, len(args))
			}
			now := float64(time.Now().UnixNano()) / 1e9
			return &object.Float{Value: now}
		},
	},
	"read_file": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", line, col, len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be string, got %s", line, col, args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("could not read file %q: %s", line, col, path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", line, col, len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to `write_file` must be (string, string), got %s and %s", line, col, args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("could not write file %q: %s", line, col, path.Value, err.Error())
			}
			return NULL
		},
	},
}
