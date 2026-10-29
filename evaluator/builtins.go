package evaluator

import (
	"fmt"
	"math"
	"os"
	"rustx/object"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
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
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", arg.Type())}
			}
		},
	},
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to str_contains. got=%d, want=2", len(args))}
			}
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to str_contains must be strings"}
			}
			return &object.Boolean{Value: strings.Contains(s.Value, substr.Value)}
		},
	},
	"str_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "str_upper expects 1 argument"}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to str_upper must be a string"}
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "str_lower expects 1 argument"}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to str_lower must be a string"}
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "math_sqrt expects 1 argument"}
			}
			val := getFloatVal(args[0])
			if val < 0 {
				return &object.Error{Message: "cannot take square root of negative number"}
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "math_pow expects 2 arguments"}
			}
			b := getFloatVal(args[0])
			p := getFloatVal(args[1])
			return &object.Float{Value: math.Pow(b, p)}
		},
	},
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "math_abs expects 1 argument"}
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
				return &object.Error{Message: "argument to math_abs must be int or float"}
			}
		},
	},
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "read_file expects 1 argument (filepath)"}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "filepath must be a string"}
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error())}
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "write_file expects 2 arguments (filepath, content)"}
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to write_file must be strings"}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func getFloatVal(obj object.Object) float64 {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value)
	case *object.Float:
		return o.Value
	default:
		return 0
	}
}
