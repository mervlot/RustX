package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
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

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VarStatement:
		var val object.Object = NULL
		if node.Value != nil {
			val = Eval(node.Value, env)
			if isError(val) {
				return val
			}
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

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		var val object.Object = NULL
		if node.ReturnValue != nil {
			val = Eval(node.ReturnValue, env)
			if isError(val) {
				return val
			}
		}
		return &object.ReturnValue{Value: val}

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
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
		return evalIndexExpression(node.GetToken(), left, index)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.GetToken(), node.Operator, right)

	case *ast.InfixExpression:
		if node.Operator == "=" {
			return evalAssignmentExpression(node, env)
		}
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
		return evalInfixExpression(node.GetToken(), node.Operator, left, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fnObj := &object.Function{Name: "", Parameters: params, Body: body, Env: env}
		if node.Name != nil {
			fnObj.Name = node.Name.Value
			env.Set(node.Name.Value, fnObj)
		}
		return fnObj

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(node.GetToken(), function, args)
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(tok token.Token, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(tok, right)
	default:
		return newError(tok, "unknown operator: %s%s", operator, right.Type())
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
		if intObj, ok := right.(*object.Integer); ok {
			if intObj.Value == 0 {
				return TRUE
			}
			return FALSE
		}
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return newError(tok, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(tok, operator, left, right)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ) &&
		(left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ):
		return evalFloatInfixExpression(tok, operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(tok, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return newError(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}
	switch l := left.(type) {
	case *object.Integer:
		r := right.(*object.Integer)
		return l.Value == r.Value
	case *object.Float:
		r := right.(*object.Float)
		return l.Value == r.Value
	case *object.Boolean:
		r := right.(*object.Boolean)
		return l.Value == r.Value
	case *object.String:
		r := right.(*object.String)
		return l.Value == r.Value
	case *object.Null:
		return true
	default:
		return left == right
	}
}

func evalIntegerInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
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
			return newError(tok, "division by zero")
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
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	var leftVal, rightVal float64

	if l, ok := left.(*object.Float); ok {
		leftVal = l.Value
	} else {
		leftVal = float64(left.(*object.Integer).Value)
	}

	if r, ok := right.(*object.Float); ok {
		rightVal = r.Value
	} else {
		rightVal = float64(right.(*object.Integer).Value)
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError(tok, "division by zero")
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
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
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
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.GetToken(), "invalid assignment target")
	}

	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}

	updatedVal, updated := env.Update(ident.Value, val)
	if !updated {
		return newError(ident.GetToken(), "identifier not found: %s", ident.Value)
	}

	return updatedVal
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.GetToken(), "postfix operator requires identifier")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(ident.GetToken(), "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		newVal := v.Value
		if node.Operator == "++" {
			newVal++
		} else if node.Operator == "--" {
			newVal--
		}
		newObj := &object.Integer{Value: newVal}
		env.Update(ident.Value, newObj)
		return v
	case *object.Float:
		newVal := v.Value
		if node.Operator == "++" {
			newVal++
		} else if node.Operator == "--" {
			newVal--
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return v
	default:
		return newError(node.GetToken(), "postfix %s operator not supported for type %s", node.Operator, val.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		scopedEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Consequence, scopedEnv)
	} else if ie.Alternative != nil {
		scopedEnv := object.NewEnclosedEnvironment(env)
		return Eval(ie.Alternative, scopedEnv)
	} else {
		return NULL
	}
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

		scopedEnv := object.NewEnclosedEnvironment(env)
		result = Eval(ws.Body, scopedEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	forEnv := object.NewEnclosedEnvironment(env)

	if fs.Init != nil {
		initRes := Eval(fs.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		bodyEnv := object.NewEnclosedEnvironment(forEnv)
		result = Eval(fs.Body, bodyEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}

		if fs.Post != nil {
			postRes := Eval(fs.Post, forEnv)
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

	return newError(node.GetToken(), "identifier not found: %s", node.Value)
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

func evalIndexExpression(tok token.Token, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(tok, left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(tok, left, index)
	default:
		return newError(tok, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(tok token.Token, array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError(tok, "index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(tok token.Token, strObj, index object.Object) object.Object {
	stringObject := strObj.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError(tok, "index out of bounds: %d", idx)
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func applyFunction(tok token.Token, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(tok, "wrong number of arguments: expected %d, got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError(tok, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Name.Value, args[paramIdx])
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
		if intObj, ok := obj.(*object.Integer); ok {
			return intObj.Value != 0
		}
		if floatObj, ok := obj.(*object.Float); ok {
			return floatObj.Value != 0
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

func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Err: errors.NewRuntimeError(tok, msg)}
}
