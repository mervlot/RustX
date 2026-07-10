package evaluator

import (
	"fmt"
	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
	"math"
	"os"
	"strings"
	"time"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node.Operator, node.Left, env, node)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != "" {
			env.Set(node.Name, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, node)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.AssignExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if !env.Update(node.Name.Value, val) {
			return newError(node, "identifier %s not found", node.Name.Value)
		}
		return val

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, node)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	default:
		return newError(node, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError(node, "unknown operator: -%s", right.Type())
}

func evalInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, node)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Cast int to float
		leftVal := float64(left.(*object.Integer).Value)
		return evalFloatInfixExpression(operator, &object.Float{Value: leftVal}, right, node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Cast int to float
		rightVal := float64(right.(*object.Integer).Value)
		return evalFloatInfixExpression(operator, left, &object.Float{Value: rightVal}, node)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right, node)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, node)
	case left.Type() != right.Type():
		return newError(node, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError(node, "division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError(node, "division by zero")
		}
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(operator string, left ast.Expression, env *object.Environment, node ast.Node) object.Object {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return newError(node, "postfix operator %s must be used on an identifier", operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(node, "identifier not found: %s", ident.Value)
	}

	if val.Type() == object.INTEGER_OBJ {
		v := val.(*object.Integer).Value
		var newVal *object.Integer
		if operator == "++" {
			newVal = &object.Integer{Value: v + 1}
		} else {
			newVal = &object.Integer{Value: v - 1}
		}
		env.Update(ident.Value, newVal)
		return val // Postfix returns original value
	}

	return newError(node, "unknown operator: %s on %s", operator, val.Type())
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node, "identifier not found: %s", node.Value)
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
	return &object.Error{Message: fmt.Sprintf(format, a...), Line: tok.Line, Column: tok.Column}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object, node ast.Node) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError(node, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			result = Eval(we.Body, env)
			if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
				return result
			}
		} else {
			break
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	childEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		Eval(fe.Init, childEnv)
	}

	var result object.Object = NULL

	for {
		if fe.Condition != nil {
			condition := Eval(fe.Condition, childEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, childEnv)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}

		if fe.Post != nil {
			Eval(fe.Post, childEnv)
		}
	}

	return result
}

func evalIndexExpression(left, index object.Object, node ast.Node) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError(node, "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

var builtins = map[string]*object.Builtin{
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
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect(), " ")
			}
			fmt.Println()
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	// math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "abs() takes exactly 1 argument"}
			}
			if args[0].Type() == object.INTEGER_OBJ {
				v := args[0].(*object.Integer).Value
				if v < 0 {
					return &object.Integer{Value: -v}
				}
				return args[0]
			}
			if args[0].Type() == object.FLOAT_OBJ {
				v := args[0].(*object.Float).Value
				return &object.Float{Value: math.Abs(v)}
			}
			return &object.Error{Message: "abs() takes int or float"}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "pow() takes exactly 2 arguments"}
			}
			v1 := 0.0
			v2 := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v1 = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v1 = args[0].(*object.Float).Value
			} else {
				return &object.Error{Message: "pow() arguments must be number"}
			}

			if args[1].Type() == object.INTEGER_OBJ {
				v2 = float64(args[1].(*object.Integer).Value)
			} else if args[1].Type() == object.FLOAT_OBJ {
				v2 = args[1].(*object.Float).Value
			} else {
				return &object.Error{Message: "pow() arguments must be number"}
			}

			return &object.Float{Value: math.Pow(v1, v2)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "sin() takes 1 argument"} }
			var v float64
			if args[0].Type() == object.INTEGER_OBJ { v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ { v = args[0].(*object.Float).Value
			} else { return &object.Error{Message: "sin() requires number"} }
			return &object.Float{Value: math.Sin(v)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "cos() takes 1 argument"} }
			var v float64
			if args[0].Type() == object.INTEGER_OBJ { v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ { v = args[0].(*object.Float).Value
			} else { return &object.Error{Message: "cos() requires number"} }
			return &object.Float{Value: math.Cos(v)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "tan() takes 1 argument"} }
			var v float64
			if args[0].Type() == object.INTEGER_OBJ { v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ { v = args[0].(*object.Float).Value
			} else { return &object.Error{Message: "tan() requires number"} }
			return &object.Float{Value: math.Tan(v)}
		},
	},
	// strings
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { return &object.Error{Message: "split() takes 2 arguments"} }
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 { return &object.Error{Message: "split() requires string arguments"} }
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts { elements[i] = &object.String{Value: p} }
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { return &object.Error{Message: "join() takes 2 arguments"} }
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 { return &object.Error{Message: "join() requires (array, string) arguments"} }
			parts := make([]string, len(arr.Elements))
			for i, e := range arr.Elements { parts[i] = e.Inspect() }
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { return &object.Error{Message: "contains() takes 2 arguments"} }
			s, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 { return &object.Error{Message: "contains() requires string arguments"} }
			if strings.Contains(s.Value, substr.Value) { return TRUE }
			return FALSE
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "to_upper() takes 1 argument"} }
			s, ok := args[0].(*object.String)
			if !ok { return &object.Error{Message: "to_upper() requires string"} }
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "to_lower() takes 1 argument"} }
			s, ok := args[0].(*object.String)
			if !ok { return &object.Error{Message: "to_lower() requires string"} }
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	// time
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "sleep() takes 1 argument (ms)"} }
			if args[0].Type() != object.INTEGER_OBJ { return &object.Error{Message: "sleep() requires integer ms"} }
			ms := args[0].(*object.Integer).Value
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return NULL
		},
	},
	// filesystem
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "read_file() takes 1 argument"} }
			path, ok := args[0].(*object.String)
			if !ok { return &object.Error{Message: "read_file() requires string"} }
			content, err := os.ReadFile(path.Value)
			if err != nil { return &object.Error{Message: err.Error()} }
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 { return &object.Error{Message: "write_file() takes 2 arguments"} }
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 { return &object.Error{Message: "write_file() requires strings"} }
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil { return &object.Error{Message: err.Error()} }
			return NULL
		},
	},
}
