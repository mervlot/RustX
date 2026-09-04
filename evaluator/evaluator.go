package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Program & Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if node.Value == nil {
			val = NULL
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ShortVarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ReturnStatement:
		if node.ReturnValue == nil {
			return &object.ReturnValue{Value: NULL}
		}
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL

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
		return evalIndexExpression(left, index, node.GetToken())

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.GetToken())

	case *ast.InfixExpression:
		// Handle short-circuiting for && and ||
		if node.Operator == "&&" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}
		if node.Operator == "||" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.GetToken())

	case *ast.AssignmentExpression:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if _, ok := env.Update(node.Name.Value, val); !ok {
			return newErrorPos(node.GetToken(), "identifier not found: %s", node.Name.Value)
		}
		return val

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != nil {
			env.Set(node.Name.Value, fn)
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
		return applyFunction(function, args, node.GetToken())
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

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

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}

		blockEnv := object.NewEnclosedEnvironment(env)
		result = evalBlockStatement(ws.Body, blockEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	if fs.Init != nil {
		initRes := Eval(fs.Init, loopEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fs.Condition != nil {
			cond := Eval(fs.Condition, loopEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}

		bodyEnv := object.NewEnclosedEnvironment(loopEnv)
		result = evalBlockStatement(fs.Body, bodyEnv)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}

		if fs.Post != nil {
			postRes := Eval(fs.Post, loopEnv)
			if isError(postRes) {
				return postRes
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

	return newErrorPos(node.GetToken(), "identifier not found: %s", node.Value)
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object, tok token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, tok)
	default:
		return newErrorPos(tok, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return newErrorPos(tok, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case (left.Type() == object.INTEGER_OBJ || left.Type() == object.FLOAT_OBJ) &&
		(right.Type() == object.INTEGER_OBJ || right.Type() == object.FLOAT_OBJ):
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newErrorPos(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
			return newErrorPos(tok, "division by zero")
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
		return newErrorPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	leftVal, _ := getFloatValue(left)
	rightVal, _ := getFloatValue(right)

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newErrorPos(tok, "division by zero")
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
		return newErrorPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
		return newErrorPos(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}
	switch l := left.(type) {
	case *object.Integer:
		return l.Value == right.(*object.Integer).Value
	case *object.Float:
		return l.Value == right.(*object.Float).Value
	case *object.Boolean:
		return l.Value == right.(*object.Boolean).Value
	case *object.String:
		return l.Value == right.(*object.String).Value
	case *object.Null:
		return true
	default:
		return left == right
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newErrorPos(node.GetToken(), "postfix operator requires identifier")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newErrorPos(node.GetToken(), "identifier not found: %s", ident.Value)
	}

	switch num := val.(type) {
	case *object.Integer:
		newVal := num.Value
		if node.Operator == "++" {
			newVal++
		} else if node.Operator == "--" {
			newVal--
		}
		updated := &object.Integer{Value: newVal}
		env.Update(ident.Value, updated)
		return num
	case *object.Float:
		newVal := num.Value
		if node.Operator == "++" {
			newVal += 1.0
		} else if node.Operator == "--" {
			newVal -= 1.0
		}
		updated := &object.Float{Value: newVal}
		env.Update(ident.Value, updated)
		return num
	default:
		return newErrorPos(node.GetToken(), "postfix operator applied to non-numeric type: %s", val.Type())
	}
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

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, tok)
	default:
		return newErrorPos(tok, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, tok token.Token) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newErrorPos(tok, "index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, tok token.Token) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newErrorPos(tok, "string index out of bounds: %d", idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
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

func applyFunction(fn object.Object, args []object.Object, tok token.Token) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newErrorPos(tok, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Name.Value, args[paramIdx])
		} else {
			env.Set(param.Name.Value, NULL)
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

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		if num, ok := obj.(*object.Integer); ok {
			return num.Value != 0
		}
		if str, ok := obj.(*object.String); ok {
			return len(str.Value) > 0
		}
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func newErrorPos(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}
