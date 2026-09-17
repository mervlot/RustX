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
	"println": {
		Fn: func(args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return &object.Null{}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
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
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to `contains` must be strings"}
			}
			return &object.Boolean{Value: strings.Contains(s.Value, sub.Value)}
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
				return &object.Error{Message: "arguments to `split` must be strings"}
			}
			parts := strings.Split(s.Value, sep.Value)
			elems := make([]object.Object, len(parts))
			for i, p := range parts {
				elems[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elems}
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
				return &object.Error{Message: "first argument to `join` must be array and second string"}
			}
			strs := []string{}
			for _, elem := range arr.Elements {
				strs = append(strs, elem.Inspect())
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to `to_upper` must be string"}
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
				return &object.Error{Message: "argument to `to_lower` must be string"}
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
				if arg.Value < 0 {
					return &object.Integer{Value: -arg.Value}
				}
				return arg
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `abs` not supported, got %s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			x, xOk := args[0].(*object.Float)
			if !xOk {
				if xInt, ok := args[0].(*object.Integer); ok {
					x = &object.Float{Value: float64(xInt.Value)}
				}
			}
			y, yOk := args[1].(*object.Float)
			if !yOk {
				if yInt, ok := args[1].(*object.Integer); ok {
					y = &object.Float{Value: float64(yInt.Value)}
				}
			}
			if x == nil || y == nil {
				return &object.Error{Message: "arguments to `pow` must be numbers"}
			}
			return &object.Float{Value: math.Pow(x.Value, y.Value)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			v, ok := args[0].(*object.Float)
			if !ok {
				if vInt, ok := args[0].(*object.Integer); ok {
					v = &object.Float{Value: float64(vInt.Value)}
				}
			}
			if v == nil {
				return &object.Error{Message: "argument to `sin` must be a number"}
			}
			return &object.Float{Value: math.Sin(v.Value)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			v, ok := args[0].(*object.Float)
			if !ok {
				if vInt, ok := args[0].(*object.Integer); ok {
					v = &object.Float{Value: float64(vInt.Value)}
				}
			}
			if v == nil {
				return &object.Error{Message: "argument to `cos` must be a number"}
			}
			return &object.Float{Value: math.Cos(v.Value)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			v, ok := args[0].(*object.Float)
			if !ok {
				if vInt, ok := args[0].(*object.Integer); ok {
					v = &object.Float{Value: float64(vInt.Value)}
				}
			}
			if v == nil {
				return &object.Error{Message: "argument to `tan` must be a number"}
			}
			return &object.Float{Value: math.Tan(v.Value)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to `read_file` must be string"}
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
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to `write_file` must be strings"}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
			}
			return &object.Boolean{Value: true}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=0", len(args))}
			}
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: "argument to `sleep` must be an integer (milliseconds)"}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
}
