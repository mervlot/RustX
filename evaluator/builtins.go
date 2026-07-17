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
			strArgs := make([]string, len(args))
			for i, arg := range args {
				strArgs[i] = arg.Inspect()
			}
			fmt.Print(strings.Join(strArgs, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			strArgs := make([]string, len(args))
			for i, arg := range args {
				strArgs[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(strArgs, " "))
			return NULL
		},
	},

	// Strings + Arrays
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("len: expected 1 argument, got %d", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: fmt.Sprintf("len: argument to 'len' not supported, got %s", args[0].Type())}
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("contains: expected 2 arguments, got %d", len(args))}
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("contains: expected strings, got (%s, %s)", args[0].Type(), args[1].Type())}
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("split: expected 2 arguments, got %d", len(args))}
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("split: expected strings, got (%s, %s)", args[0].Type(), args[1].Type())}
			}
			parts := strings.Split(str.Value, sep.Value)
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
				return &object.Error{Message: fmt.Sprintf("join: expected 2 arguments, got %d", len(args))}
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("join: expected (Array, String), got (%s, %s)", args[0].Type(), args[1].Type())}
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
				return &object.Error{Message: fmt.Sprintf("to_upper: expected 1 argument, got %d", len(args))}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("to_upper: expected string, got %s", args[0].Type())}
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("to_lower: expected 1 argument, got %d", len(args))}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("to_lower: expected string, got %s", args[0].Type())}
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},

	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("abs: expected 1 argument, got %d", len(args))}
			}
			switch val := args[0].(type) {
			case *object.Integer:
				v := val.Value
				if v < 0 {
					v = -v
				}
				return &object.Integer{Value: v}
			case *object.Float:
				return &object.Float{Value: math.Abs(val.Value)}
			default:
				return &object.Error{Message: fmt.Sprintf("abs: expected numerical type, got %s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("pow: expected 2 arguments, got %d", len(args))}
			}
			var base, exp float64
			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return &object.Error{Message: fmt.Sprintf("pow: expected numerical base, got %s", args[0].Type())}
			}
			switch e := args[1].(type) {
			case *object.Integer:
				exp = float64(e.Value)
			case *object.Float:
				exp = e.Value
			default:
				return &object.Error{Message: fmt.Sprintf("pow: expected numerical exponent, got %s", args[1].Type())}
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("sin: expected 1 argument, got %d", len(args))}
			}
			var num float64
			switch n := args[0].(type) {
			case *object.Integer:
				num = float64(n.Value)
			case *object.Float:
				num = n.Value
			default:
				return &object.Error{Message: fmt.Sprintf("sin: expected numerical type, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Sin(num)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("cos: expected 1 argument, got %d", len(args))}
			}
			var num float64
			switch n := args[0].(type) {
			case *object.Integer:
				num = float64(n.Value)
			case *object.Float:
				num = n.Value
			default:
				return &object.Error{Message: fmt.Sprintf("cos: expected numerical type, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Cos(num)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("tan: expected 1 argument, got %d", len(args))}
			}
			var num float64
			switch n := args[0].(type) {
			case *object.Integer:
				num = float64(n.Value)
			case *object.Float:
				num = n.Value
			default:
				return &object.Error{Message: fmt.Sprintf("tan: expected numerical type, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Tan(num)}
		},
	},

	// Time
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("sleep: expected 1 argument, got %d", len(args))}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("sleep: expected integer milliseconds, got %s", args[0].Type())}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("now: expected 0 arguments, got %d", len(args))}
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("read_file: expected 1 argument, got %d", len(args))}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("read_file: expected string path, got %s", args[0].Type())}
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file: %s", err.Error())}
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("write_file: expected 2 arguments, got %d", len(args))}
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: fmt.Sprintf("write_file: expected strings, got (%s, %s)", args[0].Type(), args[1].Type())}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file: %s", err.Error())}
			}
			return NULL
		},
	},
}
