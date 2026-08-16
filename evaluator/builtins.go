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

	// Strings / Collections
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
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
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
				return &object.Error{Message: "arguments to `contains` must be STRING"}
			}
			if strings.Contains(s.Value, substr.Value) {
				return TRUE
			}
			return FALSE
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
				return &object.Error{Message: "arguments to `split` must be STRING"}
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
				return &object.Error{Message: "arguments to `join` must be ARRAY and STRING"}
			}
			var strParts []string
			for _, elem := range arr.Elements {
				strParts = append(strParts, elem.Inspect())
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
				return &object.Error{Message: "argument to `to_upper` must be STRING"}
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
				return &object.Error{Message: "argument to `to_lower` must be STRING"}
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: "argument to `sleep` must be INTEGER"}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
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
				return &object.Error{Message: fmt.Sprintf("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			var x, y float64
			if i, ok := args[0].(*object.Integer); ok {
				x = float64(i.Value)
			} else if f, ok := args[0].(*object.Float); ok {
				x = f.Value
			} else {
				return &object.Error{Message: fmt.Sprintf("argument 1 to `pow` must be numeric, got %s", args[0].Type())}
			}

			if i, ok := args[1].(*object.Integer); ok {
				y = float64(i.Value)
			} else if f, ok := args[1].(*object.Float); ok {
				y = f.Value
			} else {
				return &object.Error{Message: fmt.Sprintf("argument 2 to `pow` must be numeric, got %s", args[1].Type())}
			}

			return &object.Float{Value: math.Pow(x, y)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var x float64
			if i, ok := args[0].(*object.Integer); ok {
				x = float64(i.Value)
			} else if f, ok := args[0].(*object.Float); ok {
				x = f.Value
			} else {
				return &object.Error{Message: fmt.Sprintf("argument to `sin` must be numeric, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Sin(x)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var x float64
			if i, ok := args[0].(*object.Integer); ok {
				x = float64(i.Value)
			} else if f, ok := args[0].(*object.Float); ok {
				x = f.Value
			} else {
				return &object.Error{Message: fmt.Sprintf("argument to `cos` must be numeric, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Cos(x)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			var x float64
			if i, ok := args[0].(*object.Integer); ok {
				x = float64(i.Value)
			} else if f, ok := args[0].(*object.Float); ok {
				x = f.Value
			} else {
				return &object.Error{Message: fmt.Sprintf("argument to `tan` must be numeric, got %s", args[0].Type())}
			}
			return &object.Float{Value: math.Tan(x)}
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument to `read_file` must be STRING, got %s", args[0].Type())}
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error())}
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument 1 to `write_file` must be STRING, got %s", args[0].Type())}
			}
			content, ok := args[1].(*object.String)
			if !ok {
				return &object.Error{Message: fmt.Sprintf("argument 2 to `write_file` must be STRING, got %s", args[1].Type())}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
			}
			return TRUE
		},
	},
}
