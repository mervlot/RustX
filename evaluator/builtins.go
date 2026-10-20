package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				if arg != nil {
					out = append(out, arg.Inspect())
				}
			}
			fmt.Print(strings.Join(out, " "))
			return &object.Null{}
		},
	},
	"println": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			out := []string{}
			for _, arg := range args {
				if arg != nil {
					out = append(out, arg.Inspect())
				}
			}
			fmt.Println(strings.Join(out, " "))
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type()), Pos: pos}
			}
		},
	},
	"split": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to split. got=%d, want=2", len(args)), Pos: pos}
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to split must be strings", Pos: pos}
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
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to join. got=%d, want=2", len(args)), Pos: pos}
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "first argument to join must be array, second must be string", Pos: pos}
			}
			strs := []string{}
			for _, elem := range arr.Elements {
				strs = append(strs, elem.Inspect())
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to to_upper must be string", Pos: pos}
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to to_lower must be string", Pos: pos}
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"contains": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to contains. got=%d, want=2", len(args)), Pos: pos}
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to contains must be strings", Pos: pos}
			}
			return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
		},
	},
	"abs": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to abs. got=%d, want=1", len(args)), Pos: pos}
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
				return &object.Error{Message: fmt.Sprintf("argument to abs not supported, got %s", args[0].Type()), Pos: pos}
			}
		},
	},
	"pow": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to pow. got=%d, want=2", len(args)), Pos: pos}
			}
			b, ok1 := floatValue(args[0])
			p, ok2 := floatValue(args[1])
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to pow must be numbers", Pos: pos}
			}
			res := math.Pow(b, p)
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			val, ok := floatValue(args[0])
			if !ok {
				return &object.Error{Message: "argument to sin must be number", Pos: pos}
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			val, ok := floatValue(args[0])
			if !ok {
				return &object.Error{Message: "argument to cos must be number", Pos: pos}
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args)), Pos: pos}
			}
			val, ok := floatValue(args[0])
			if !ok {
				return &object.Error{Message: "argument to tan must be number", Pos: pos}
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to sleep. got=%d, want=1", len(args)), Pos: pos}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: "argument to sleep must be integer (milliseconds)", Pos: pos}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"read_file": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to read_file. got=%d, want=1", len(args)), Pos: pos}
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to read_file must be string", Pos: pos}
			}
			content, err := os.ReadFile(filepath.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error()), Pos: pos}
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(pos token.Position, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to write_file. got=%d, want=2", len(args)), Pos: pos}
			}
			filepath, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "arguments to write_file must be strings", Pos: pos}
			}
			err := os.WriteFile(filepath.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error()), Pos: pos}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func floatValue(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
