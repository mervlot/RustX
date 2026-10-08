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
			out := make([]string, len(args))
			for i, arg := range args {
				out[i] = arg.Inspect()
			}
			fmt.Println(strings.Join(out, " "))
			return &object.Null{}
		},
	},
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
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to contains. got=%d, want=2", len(args))}
			}
			s, sOk := args[0].(*object.String)
			substr, subOk := args[1].(*object.String)
			if !sOk || !subOk {
				return &object.Error{Message: "arguments to contains must be strings"}
			}
			return &object.Boolean{Value: strings.Contains(s.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to split. got=%d, want=2", len(args))}
			}
			s, sOk := args[0].(*object.String)
			sep, sepOk := args[1].(*object.String)
			if !sOk || !sepOk {
				return &object.Error{Message: "arguments to split must be strings"}
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
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to join. got=%d, want=2", len(args))}
			}
			arr, arrOk := args[0].(*object.Array)
			sep, sepOk := args[1].(*object.String)
			if !arrOk || !sepOk {
				return &object.Error{Message: "first argument to join must be an array and second a string"}
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
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to to_upper. got=%d, want=1", len(args))}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to to_upper must be a string"}
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to to_lower. got=%d, want=1", len(args))}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to to_lower must be a string"}
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to abs. got=%d, want=1", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				val := arg.Value
				if val < 0 {
					val = -val
				}
				return &object.Integer{Value: val}
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to abs not supported, got %s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to pow. got=%d, want=2", len(args))}
			}
			base, bOk := getFloatVal(args[0])
			exp, eOk := getFloatVal(args[1])
			if !bOk || !eOk {
				return &object.Error{Message: "arguments to pow must be numbers"}
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to sin. got=%d, want=1", len(args))}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "argument to sin must be a number"}
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to cos. got=%d, want=1", len(args))}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "argument to cos must be a number"}
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to tan. got=%d, want=1", len(args))}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "argument to tan must be a number"}
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to sleep. got=%d, want=1", len(args))}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Message: "argument to sleep must be integer milliseconds"}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to now. got=%d, want=0", len(args))}
			}
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to read_file. got=%d, want=1", len(args))}
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "argument to read_file must be a string path"}
			}
			data, err := os.ReadFile(filepath.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error())}
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments to write_file. got=%d, want=2", len(args))}
			}
			filepath, fOk := args[0].(*object.String)
			content, cOk := args[1].(*object.String)
			if !fOk || !cOk {
				return &object.Error{Message: "arguments to write_file must be (path, content)"}
			}
			err := os.WriteFile(filepath.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func getFloatVal(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
