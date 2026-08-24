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
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
			return &object.Null{}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return &object.Null{}
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("len() takes exactly 1 argument (%d given)", len(args))}
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Err: fmt.Errorf("argument to len() not supported, got %s", args[0].Type())}
			}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Err: fmt.Errorf("contains() takes 2 arguments (%d given)", len(args))}
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Err: fmt.Errorf("arguments to contains() must be strings")}
			}
			return &object.Boolean{Value: strings.Contains(str.Value, substr.Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Err: fmt.Errorf("split() takes 2 arguments (%d given)", len(args))}
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Err: fmt.Errorf("arguments to split() must be strings")}
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
				return &object.Error{Err: fmt.Errorf("join() takes 2 arguments (%d given)", len(args))}
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Err: fmt.Errorf("arguments to join() must be (Array, String)")}
			}
			strs := make([]string, len(arr.Elements))
			for i, elem := range arr.Elements {
				strs[i] = elem.Inspect()
			}
			return &object.String{Value: strings.Join(strs, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("to_upper() takes 1 argument")}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to to_upper() must be String")}
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("to_lower() takes 1 argument")}
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to to_lower() must be String")}
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("abs() takes 1 argument")}
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
				return &object.Error{Err: fmt.Errorf("argument to abs() must be number, got %s", args[0].Type())}
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Err: fmt.Errorf("pow() takes 2 arguments")}
			}
			base, ok1 := getFloatValue(args[0])
			exp, ok2 := getFloatValue(args[1])
			if !ok1 || !ok2 {
				return &object.Error{Err: fmt.Errorf("arguments to pow() must be numbers")}
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("sin() takes 1 argument")}
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to sin() must be number")}
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("cos() takes 1 argument")}
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to cos() must be number")}
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("tan() takes 1 argument")}
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to tan() must be number")}
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("sleep() takes 1 argument (milliseconds)")}
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to sleep() must be integer milliseconds")}
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return &object.Null{}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano()}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Err: fmt.Errorf("read_file() takes 1 argument (filepath)")}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Err: fmt.Errorf("argument to read_file() must be string")}
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Err: fmt.Errorf("read_file failed: %v", err)}
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Err: fmt.Errorf("write_file() takes 2 arguments (filepath, content)")}
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Err: fmt.Errorf("arguments to write_file() must be strings")}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Err: fmt.Errorf("write_file failed: %v", err)}
			}
			return &object.Boolean{Value: true}
		},
	},
}

func getFloatValue(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
