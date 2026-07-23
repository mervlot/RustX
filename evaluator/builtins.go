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
			strArgs := make([]string, len(args))
			for i, arg := range args {
				strArgs[i] = arg.Inspect()
			}
			fmt.Print(strings.Join(strArgs, " "))
			return &object.Null{}
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			strArgs := make([]string, len(args))
			for i, arg := range args {
				strArgs[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(strArgs, " "))
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to 'len' not supported, got %s", args[0].Type())}
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("arguments to 'contains' must be strings. got=%s, %s", args[0].Type(), args[1].Type())}
			}
			return &object.Boolean{Value: strings.Contains(s.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("arguments to 'split' must be strings. got=%s, %s", args[0].Type(), args[1].Type())}
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
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("arguments to 'join' must be array and string. got=%s, %s", args[0].Type(), args[1].Type())}
			}
			strParts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				strParts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(strParts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to 'to_upper' must be string. got=%s", args[0].Type())}
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to 'to_lower' must be string. got=%s", args[0].Type())}
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
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
				return &object.Error{Message: fmt.Sprintf("argument to 'abs' must be integer or float. got=%s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			var base, exp float64
			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return &object.Error{Message: fmt.Sprintf("first argument to 'pow' must be numeric. got=%s", args[0].Type())}
			}
			switch e := args[1].(type) {
			case *object.Integer:
				exp = float64(e.Value)
			case *object.Float:
				exp = e.Value
			default:
				return &object.Error{Message: fmt.Sprintf("second argument to 'pow' must be numeric. got=%s", args[1].Type())}
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return &object.Error{Message: fmt.Sprintf("argument to 'sin' must be numeric. got=%s", args[0].Type())}
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return &object.Error{Message: fmt.Sprintf("argument to 'cos' must be numeric. got=%s", args[0].Type())}
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return &object.Error{Message: fmt.Sprintf("argument to 'tan' must be numeric. got=%s", args[0].Type())}
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to 'sleep' must be integer milliseconds. got=%s", args[0].Type())}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=0", len(args))}
			}
			return &object.Integer{Value: time.Now().UnixMilli()}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to 'read_file' must be string. got=%s", args[0].Type())}
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("failed to read file: %s", err.Error())}
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("arguments to 'write_file' must be strings. got=%s, %s", args[0].Type(), args[1].Type())}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("failed to write file: %s", err.Error())}
			}
			return &object.Null{}
		},
	},
}
