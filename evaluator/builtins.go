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
	// Printing
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Print(strings.Join(parts, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(strings.Join(parts, " "))
			return NULL
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
	"type": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},

	// Math functions
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("abs() takes exactly 1 argument")
			}
			switch v := args[0].(type) {
			case *object.Integer:
				if v.Value < 0 {
					return &object.Integer{Value: -v.Value}
				}
				return v
			case *object.Float:
				return &object.Float{Value: math.Abs(v.Value)}
			default:
				return newError("abs() argument must be integer or float, got %s", args[0].Type())
			}
		},
	},
	"sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("sqrt() takes exactly 1 argument")
			}
			f := getFloat(args[0])
			if f == nil {
				return newError("sqrt() requires a number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sqrt(*f)}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("pow() takes exactly 2 arguments")
			}
			x := getFloat(args[0])
			y := getFloat(args[1])
			if x == nil || y == nil {
				return newError("pow() requires two numeric arguments")
			}
			return &object.Float{Value: math.Pow(*x, *y)}
		},
	},
	"floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("floor() takes exactly 1 argument")
			}
			f := getFloat(args[0])
			if f == nil {
				return newError("floor() requires a number")
			}
			return &object.Float{Value: math.Floor(*f)}
		},
	},
	"ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("ceil() takes exactly 1 argument")
			}
			f := getFloat(args[0])
			if f == nil {
				return newError("ceil() requires a number")
			}
			return &object.Float{Value: math.Ceil(*f)}
		},
	},
	"min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("min() requires at least 2 arguments")
			}
			minVal := getFloat(args[0])
			if minVal == nil {
				return newError("min() requires numeric arguments")
			}
			allInt := args[0].Type() == object.INTEGER_OBJ
			res := *minVal
			for _, a := range args[1:] {
				f := getFloat(a)
				if f == nil {
					return newError("min() requires numeric arguments")
				}
				if a.Type() != object.INTEGER_OBJ {
					allInt = false
				}
				if *f < res {
					res = *f
				}
			}
			if allInt {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("max() requires at least 2 arguments")
			}
			maxVal := getFloat(args[0])
			if maxVal == nil {
				return newError("max() requires numeric arguments")
			}
			allInt := args[0].Type() == object.INTEGER_OBJ
			res := *maxVal
			for _, a := range args[1:] {
				f := getFloat(a)
				if f == nil {
					return newError("max() requires numeric arguments")
				}
				if a.Type() != object.INTEGER_OBJ {
					allInt = false
				}
				if *f > res {
					res = *f
				}
			}
			if allInt {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},

	// Strings
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("to_upper() takes 1 string argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("to_upper() requires a string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("to_lower() takes 1 string argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("to_lower() requires a string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("trim() takes 1 string argument")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("trim() requires a string")
			}
			return &object.String{Value: strings.TrimSpace(s.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("contains() takes 2 arguments (str, substr)")
			}
			s1, ok1 := args[0].(*object.String)
			s2, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("contains() requires two string arguments")
			}
			return nativeBoolToBooleanObject(strings.Contains(s1.Value, s2.Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("split() takes 2 arguments (str, sep)")
			}
			s1, ok1 := args[0].(*object.String)
			s2, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("split() requires two string arguments")
			}
			parts := strings.Split(s1.Value, s2.Value)
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
				return newError("join() takes 2 arguments (array, sep)")
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("join() requires an array and a string separator")
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},

	// Time
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("sleep() takes 1 argument (milliseconds)")
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("sleep() argument must be an integer (milliseconds)")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("read_file() takes 1 argument (filepath)")
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError("read_file() argument must be a string")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError("read_file error: %s", err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("write_file() takes 2 arguments (filepath, content)")
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("write_file() requires string filepath and content")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("write_file error: %s", err.Error())
			}
			return TRUE
		},
	},
}

func getFloat(obj object.Object) *float64 {
	switch v := obj.(type) {
	case *object.Integer:
		f := float64(v.Value)
		return &f
	case *object.Float:
		f := v.Value
		return &f
	default:
		return nil
	}
}
