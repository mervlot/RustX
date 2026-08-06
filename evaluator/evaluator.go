package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(line int, col int, format string, args ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Column:  col,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
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
		// Any non-null/non-false is truthy
		return true
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	if node == nil {
		return NULL
	}

	switch node := node.(type) {
	// Program
	case *ast.Program:
		return evalProgram(node, env)

	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ShorthandStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		success := env.Update(node.Name.Value, val)
		if !success {
			return newError(node.Token.Line, node.Token.Column, "identifier not found: %s", node.Name.Value)
		}
		return NULL

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

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

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		if node.Operator == "&&" {
			leftVal := Eval(node.Left, env)
			if isError(leftVal) {
				return leftVal
			}
			if !isTruthy(leftVal) {
				return nativeBoolToBooleanObject(false)
			}
			rightVal := Eval(node.Right, env)
			if isError(rightVal) {
				return rightVal
			}
			return nativeBoolToBooleanObject(isTruthy(rightVal))
		}
		if node.Operator == "||" {
			leftVal := Eval(node.Left, env)
			if isError(leftVal) {
				return leftVal
			}
			if isTruthy(leftVal) {
				return nativeBoolToBooleanObject(true)
			}
			rightVal := Eval(node.Right, env)
			if isError(rightVal) {
				return rightVal
			}
			return nativeBoolToBooleanObject(isTruthy(rightVal))
		}

		leftVal := Eval(node.Left, env)
		if isError(leftVal) {
			return leftVal
		}
		rightVal := Eval(node.Right, env)
		if isError(rightVal) {
			return rightVal
		}
		return evalInfixExpression(node.Operator, leftVal, rightVal, node)

	case *ast.IfExpression:
		cond := Eval(node.Condition, env)
		if isError(cond) {
			return cond
		}
		if isTruthy(cond) {
			return Eval(node.Consequence, env)
		} else if node.Alternative != nil {
			return Eval(node.Alternative, env)
		}
		return NULL

	case *ast.WhileExpression:
		var lastVal object.Object = NULL
		for {
			cond := Eval(node.Condition, env)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
			lastVal = Eval(node.Body, env)
			if isError(lastVal) {
				return lastVal
			}
			if lastVal.Type() == object.RETURN_VALUE_OBJ {
				return lastVal
			}
		}
		return lastVal

	case *ast.ForExpression:
		forEnv := object.NewEnclosedEnvironment(env)

		if node.Init != nil {
			initResult := Eval(node.Init, forEnv)
			if isError(initResult) {
				return initResult
			}
		}

		var lastVal object.Object = NULL
		for {
			if node.Condition != nil {
				cond := Eval(node.Condition, forEnv)
				if isError(cond) {
					return cond
				}
				if !isTruthy(cond) {
					break
				}
			}

			lastVal = Eval(node.Body, forEnv)
			if isError(lastVal) {
				return lastVal
			}
			if lastVal.Type() == object.RETURN_VALUE_OBJ {
				return lastVal
			}

			if node.Post != nil {
				postResult := Eval(node.Post, forEnv)
				if isError(postResult) {
					return postResult
				}
			}
		}
		return lastVal

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
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
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch val := result.(type) {
		case *object.ReturnValue:
			return val.Value
		case *object.Error:
			return val
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Token.Line, node.Token.Column, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(op string, right object.Object, node ast.Expression) object.Object {
	switch op {
	case "!":
		return nativeBoolToBooleanObject(!isTruthy(right))
	case "-":
		switch right := right.(type) {
		case *object.Integer:
			return &object.Integer{Value: -right.Value}
		case *object.Float:
			return &object.Float{Value: -right.Value}
		default:
			return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: -%s", right.Type())
		}
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: %s%s", op, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.Token.Line, node.Token.Column, "operator %s can only be applied to identifiers", node.Operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(ident.Token.Line, ident.Token.Column, "identifier not found: %s", ident.Value)
	}

	var originalVal object.Object
	var newVal object.Object

	switch val := val.(type) {
	case *object.Integer:
		originalVal = &object.Integer{Value: val.Value}
		var delta int64 = 1
		if node.Operator == "--" {
			delta = -1
		}
		newVal = &object.Integer{Value: val.Value + delta}

	case *object.Float:
		originalVal = &object.Float{Value: val.Value}
		var delta float64 = 1.0
		if node.Operator == "--" {
			delta = -1.0
		}
		newVal = &object.Float{Value: val.Value + delta}

	default:
		return newError(node.Token.Line, node.Token.Column, "operator %s requires integer or float type, got %s", node.Operator, val.Type())
	}

	env.Update(ident.Value, newVal)
	return originalVal
}

func evalInfixExpression(op string, left, right object.Object, node ast.Expression) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left.(*object.Integer), right.(*object.Integer), node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(op, left.(*object.Float), right.(*object.Float), node)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		lFloat := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(op, lFloat, right.(*object.Float), node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		rFloat := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(op, left.(*object.Float), rFloat, node)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left.(*object.String), right.(*object.String), node)
	case op == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case op == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right *object.Integer, node ast.Expression) object.Object {
	lVal := left.Value
	rVal := right.Value

	switch op {
	case "+":
		return &object.Integer{Value: lVal + rVal}
	case "-":
		return &object.Integer{Value: lVal - rVal}
	case "*":
		return &object.Integer{Value: lVal * rVal}
	case "/":
		if rVal == 0 {
			return newError(node.GetToken().Line, node.GetToken().Column, "division by zero")
		}
		return &object.Integer{Value: lVal / rVal}
	case "<":
		return nativeBoolToBooleanObject(lVal < rVal)
	case ">":
		return nativeBoolToBooleanObject(lVal > rVal)
	case "<=":
		return nativeBoolToBooleanObject(lVal <= rVal)
	case ">=":
		return nativeBoolToBooleanObject(lVal >= rVal)
	case "==":
		return nativeBoolToBooleanObject(lVal == rVal)
	case "!=":
		return nativeBoolToBooleanObject(lVal != rVal)
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalFloatInfixExpression(op string, left, right *object.Float, node ast.Expression) object.Object {
	lVal := left.Value
	rVal := right.Value

	switch op {
	case "+":
		return &object.Float{Value: lVal + rVal}
	case "-":
		return &object.Float{Value: lVal - rVal}
	case "*":
		return &object.Float{Value: lVal * rVal}
	case "/":
		if rVal == 0.0 {
			return newError(node.GetToken().Line, node.GetToken().Column, "division by zero")
		}
		return &object.Float{Value: lVal / rVal}
	case "<":
		return nativeBoolToBooleanObject(lVal < rVal)
	case ">":
		return nativeBoolToBooleanObject(lVal > rVal)
	case "<=":
		return nativeBoolToBooleanObject(lVal <= rVal)
	case ">=":
		return nativeBoolToBooleanObject(lVal >= rVal)
	case "==":
		return nativeBoolToBooleanObject(lVal == rVal)
	case "!=":
		return nativeBoolToBooleanObject(lVal != rVal)
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right *object.String, node ast.Expression) object.Object {
	lVal := left.Value
	rVal := right.Value

	switch op {
	case "+":
		return &object.String{Value: lVal + rVal}
	case "==":
		return nativeBoolToBooleanObject(lVal == rVal)
	case "!=":
		return nativeBoolToBooleanObject(lVal != rVal)
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
			return float64(left.(*object.Integer).Value) == right.(*object.Float).Value
		}
		if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
			return left.(*object.Float).Value == float64(right.(*object.Integer).Value)
		}
		return false
	}

	switch left.Type() {
	case object.INTEGER_OBJ:
		return left.(*object.Integer).Value == right.(*object.Integer).Value
	case object.FLOAT_OBJ:
		return left.(*object.Float).Value == right.(*object.Float).Value
	case object.BOOLEAN_OBJ:
		return left.(*object.Boolean).Value == right.(*object.Boolean).Value
	case object.STRING_OBJ:
		return left.(*object.String).Value == right.(*object.String).Value
	case object.NULL_OBJ:
		return true
	default:
		return false
	}
}

func evalIndexExpression(left, index object.Object, node ast.Expression) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left.(*object.Array), index.(*object.Integer), node)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left.(*object.String), index.(*object.Integer), node)
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer, node ast.Expression) object.Object {
	idx := index.Value
	maxIdx := int64(len(array.Elements) - 1)

	if idx < 0 || idx > maxIdx {
		return newError(node.GetToken().Line, node.GetToken().Column, "index out of bounds: %d", idx)
	}

	return array.Elements[idx]
}

func evalStringIndexExpression(str *object.String, index *object.Integer, node ast.Expression) object.Object {
	idx := index.Value
	maxIdx := int64(len(str.Value) - 1)

	if idx < 0 || idx > maxIdx {
		return newError(node.GetToken().Line, node.GetToken().Column, "index out of bounds: %d", idx)
	}

	return &object.String{Value: string(str.Value[idx])}
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

func applyFunction(fn object.Object, args []object.Object, node ast.Expression) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		res := fn.Fn(args...)
		if err, ok := res.(*object.Error); ok {
			if err.Line == 0 && err.Column == 0 {
				err.Line = node.GetToken().Line
				err.Column = node.GetToken().Column
			}
		}
		return res
	default:
		return newError(node.GetToken().Line, node.GetToken().Column, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Name.Value, args[paramIdx])
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
