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

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

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

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.GetToken())

	case *ast.InfixExpression:
		// Short-circuiting evaluation for logical operations
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

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Body: body, Env: env}
		// Named function definitions are supported by registering in environment
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

	case *ast.ArrayLiteral:
		elems := evalExpressions(node.Elements, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

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
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // Initialize result to prevention of nil pointers

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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorAtToken(node.GetToken(), "identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object, tok token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, tok)
	default:
		return newErrorAtToken(tok, "unknown operator: %s%s", operator, right.Type())
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
		return newErrorAtToken(tok, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer), tok)
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String), tok)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newErrorAtToken(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	// Compare value types properly to ensure constant value evaluation
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

func evalIntegerInfixExpression(operator string, left, right *object.Integer, tok token.Token) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newErrorAtToken(tok, "division by zero")
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
		return newErrorAtToken(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	lVal, ok1 := getFloatValue(left)
	rVal, ok2 := getFloatValue(right)
	if !ok1 || !ok2 {
		return newErrorAtToken(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	switch operator {
	case "+":
		return &object.Float{Value: lVal + rVal}
	case "-":
		return &object.Float{Value: lVal - rVal}
	case "*":
		return &object.Float{Value: lVal * rVal}
	case "/":
		if rVal == 0 {
			return newErrorAtToken(tok, "division by zero")
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
		return newErrorAtToken(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right *object.String, tok token.Token) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return newErrorAtToken(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newErrorAtToken(node.GetToken(), "postfix operator %s must be applied to identifier", node.Operator)
	}

	val, found := env.Get(ident.Value)
	if !found {
		return newErrorAtToken(node.GetToken(), "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		var newVal int64
		if node.Operator == "++" {
			newVal = v.Value + 1
		} else {
			newVal = v.Value - 1
		}
		newObj := &object.Integer{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj
	case *object.Float:
		var newVal float64
		if node.Operator == "++" {
			newVal = v.Value + 1
		} else {
			newVal = v.Value - 1
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj
	default:
		return newErrorAtToken(node.GetToken(), "postfix operator %s cannot be applied to type %s", node.Operator, val.Type())
	}
}

func evalAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	// Shorthand := assignment or = update
	if node.Operator == ":=" {
		ident, ok := node.Left.(*ast.Identifier)
		if !ok {
			return newErrorAtToken(node.GetToken(), "LHS of := must be identifier")
		}
		env.Set(ident.Value, val)
		return val
	}

	// Normal variable assignment update
	switch left := node.Left.(type) {
	case *ast.Identifier:
		if ok := env.Update(left.Value, val); !ok {
			// Shorthand initialization fallback or return error?
			// The spec says 'var' is for mutable variables, and ':=' is shorthand.
			// Normal '=' is update, so if it doesn't exist, we should return an error.
			return newErrorAtToken(node.GetToken(), "identifier not found for assignment: %s", left.Value)
		}
		return val
	case *ast.IndexExpression:
		// Array index assignment: arr[idx] = val
		arrObj := Eval(left.Left, env)
		if isError(arrObj) {
			return arrObj
		}
		arr, ok := arrObj.(*object.Array)
		if !ok {
			return newErrorAtToken(left.GetToken(), "cannot assign to index of type %s", arrObj.Type())
		}
		idxObj := Eval(left.Index, env)
		if isError(idxObj) {
			return idxObj
		}
		idx, ok := idxObj.(*object.Integer)
		if !ok {
			return newErrorAtToken(left.GetToken(), "array index must be integer, got %s", idxObj.Type())
		}
		indexVal := idx.Value
		if indexVal < 0 || indexVal >= int64(len(arr.Elements)) {
			return newErrorAtToken(left.GetToken(), "index out of bounds: %d", indexVal)
		}
		arr.Elements[indexVal] = val
		return val
	default:
		return newErrorAtToken(node.GetToken(), "invalid target for assignment")
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

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		result = Eval(we.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// For loop scoping: encapsulate initializer variables
	loopEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		initRes := Eval(fe.Init, loopEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fe.Condition != nil {
			condition := Eval(fe.Condition, loopEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, loopEnv)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

		if fe.Post != nil {
			postRes := Eval(fe.Post, loopEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return result
}

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left.(*object.Array), index.(*object.Integer), tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left.(*object.String), index.(*object.Integer), tok)
	default:
		return newErrorAtToken(tok, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer, tok token.Token) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return newErrorAtToken(tok, "index out of bounds: %d", idx)
	}

	return array.Elements[idx]
}

func evalStringIndexExpression(str *object.String, index *object.Integer, tok token.Token) object.Object {
	idx := index.Value
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return newErrorAtToken(tok, "index out of bounds: %d", idx)
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

func applyFunction(fn object.Object, args []object.Object, tok token.Token) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newErrorAtToken(tok, "wrong number of arguments. got=%d, want=%d", len(args), len(fn.Parameters))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newErrorAtToken(tok, "not a function: %s", fn.Type())
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

func newErrorAtToken(tok token.Token, format string, a ...interface{}) object.Object {
	msg := fmt.Sprintf(format, a...)
	return &object.Error{
		Message: fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", msg, tok.Line, tok.Column),
	}
}
