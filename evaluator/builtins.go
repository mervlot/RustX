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
		Fn: func(line int, col int, args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"println": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
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
	"abs": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Integer{Value: val}
			case *object.Float:
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Float{Value: val}
			default:
				return newError(line, col, "argument to 'abs' not supported, got=%s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, col, "wrong number of arguments. got=%d, want=2", len(args))
			}
			var base, power float64
			switch arg := args[0].(type) {
			case *object.Integer:
				base = float64(arg.Value)
			case *object.Float:
				base = arg.Value
			default:
				return newError(line, col, "first argument to 'pow' must be number, got=%s", args[0].Type())
			}
			switch arg := args[1].(type) {
			case *object.Integer:
				power = float64(arg.Value)
			case *object.Float:
				power = arg.Value
			default:
				return newError(line, col, "second argument to 'pow' must be number, got=%s", args[1].Type())
			}
			return &object.Float{Value: math.Pow(base, power)}
		},
	},
	"sin": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError(line, col, "argument to 'sin' must be number, got=%s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError(line, col, "argument to 'cos' must be number, got=%s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError(line, col, "argument to 'tan' must be number, got=%s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"len": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(line, col, "argument to 'len' not supported, got=%s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, col, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(line, col, "arguments to 'contains' must be strings")
			}
			return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, col, "wrong number of arguments. got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(line, col, "arguments to 'split' must be strings")
			}
			parts := strings.Split(str.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, part := range parts {
				elements[i] = &object.String{Value: part}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, col, "wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(line, col, "first argument to 'join' must be an array, second must be a string")
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(line, col, "argument to 'to_upper' must be string")
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(line, col, "argument to 'to_lower' must be string")
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"sleep": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(line, col, "argument to 'sleep' must be integer milliseconds")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError(line, col, "wrong number of arguments. got=%d, want=0", len(args))
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"read_file": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(line, col, "wrong number of arguments. got=%d, want=1", len(args))
			}
			pathObj, ok := args[0].(*object.String)
			if !ok {
				return newError(line, col, "argument to 'read_file' must be a string filepath")
			}
			content, err := os.ReadFile(pathObj.Value)
			if err != nil {
				return newError(line, col, "failed to read file %q: %s", pathObj.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(line int, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(line, col, "wrong number of arguments. got=%d, want=2", len(args))
			}
			pathObj, ok1 := args[0].(*object.String)
			contentObj, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(line, col, "arguments to 'write_file' must be string filepath and string content")
			}
			err := os.WriteFile(pathObj.Value, []byte(contentObj.Value), 0644)
			if err != nil {
				return &object.Boolean{Value: false}
			}
			return &object.Boolean{Value: true}
		},
	},
}
