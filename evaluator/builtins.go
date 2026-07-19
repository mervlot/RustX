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
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("first argument to `contains` must be string, got %s", args[0].Type())
			}
			substr, ok := args[1].(*object.String)
			if !ok {
				return newError("second argument to `contains` must be string, got %s", args[1].Type())
			}
			return &object.Boolean{Value: strings.Contains(s.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("first argument to `split` must be string, got %s", args[0].Type())
			}
			sep, ok := args[1].(*object.String)
			if !ok {
				return newError("second argument to `split` must be string, got %s", args[1].Type())
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
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to `join` must be array, got %s", args[0].Type())
			}
			sep, ok := args[1].(*object.String)
			if !ok {
				return newError("second argument to `join` must be string, got %s", args[1].Type())
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
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_upper` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `to_lower` must be string, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch val := args[0].(type) {
			case *object.Integer:
				if val.Value < 0 {
					return &object.Integer{Value: -val.Value}
				}
				return val
			case *object.Float:
				return &object.Float{Value: math.Abs(val.Value)}
			default:
				return newError("argument to `abs` must be integer or float, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, exp float64

			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return newError("first argument to `pow` must be integer or float, got %s", args[0].Type())
			}

			switch e := args[1].(type) {
			case *object.Integer:
				exp = float64(e.Value)
			case *object.Float:
				exp = e.Value
			default:
				return newError("second argument to `pow` must be integer or float, got %s", args[1].Type())
			}

			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `sin` must be integer or float, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `cos` must be integer or float, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("argument to `tan` must be integer or float, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Float{Value: float64(time.Now().UnixNano()) / 1e9}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var sec float64
			switch arg := args[0].(type) {
			case *object.Integer:
				sec = float64(arg.Value)
			case *object.Float:
				sec = arg.Value
			default:
				return newError("argument to `sleep` must be integer or float, got %s", args[0].Type())
			}
			time.Sleep(time.Duration(sec * float64(time.Second)))
			return &object.Null{}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `read_file` must be string, got %s", args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("failed to read file %s: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("first argument to `write_file` must be string, got %s", args[0].Type())
			}
			content, ok := args[1].(*object.String)
			if !ok {
				return newError("second argument to `write_file` must be string, got %s", args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("failed to write file %s: %s", path.Value, err.Error())
			}
			return &object.Boolean{Value: true}
		},
	},
}
