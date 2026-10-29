package evaluator

import (
	"fmt"
	"rustx/ast"
	"rustx/object"
	"rustx/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

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

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		switch target := node.Name.(type) {
		case *ast.Identifier:
			updated := env.Update(target.Value, val)
			if !updated {
				// If not previously defined, set in current env
				env.Set(target.Value, val)
			}
			return val
		case *ast.IndexExpression:
			arrayObj := Eval(target.Left, env)
			if isError(arrayObj) {
				return arrayObj
			}
			indexObj := Eval(target.Index, env)
			if isError(indexObj) {
				return indexObj
			}
			arr, ok := arrayObj.(*object.Array)
			if !ok {
				return newError(target.Left.GetToken(), "index operator not supported on %s", arrayObj.Type())
			}
			idx, ok := indexObj.(*object.Integer)
			if !ok {
				return newError(target.Index.GetToken(), "index must be integer")
			}
			if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
				return newError(target.Index.GetToken(), "index out of bounds: %d", idx.Value)
			}
			arr.Elements[idx.Value] = val
			return val
		default:
			return newError(node.GetToken(), "invalid assignment target")
		}

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

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

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.GetToken(), node.Operator, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.GetToken(), node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Name:       "",
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		if node.Name != nil {
			fn.Name = node.Name.Value
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
		return applyFunction(node.GetToken(), function, args)

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
	}

	return NULL
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

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.GetToken(), "postfix operator %s requires identifier", node.Operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(ident.GetToken(), "identifier not found: %s", ident.Value)
	}

	switch node.Operator {
	case "++":
		switch v := val.(type) {
		case *object.Integer:
			newVal := &object.Integer{Value: v.Value + 1}
			env.Update(ident.Value, newVal)
			return v
		case *object.Float:
			newVal := &object.Float{Value: v.Value + 1.0}
			env.Update(ident.Value, newVal)
			return v
		default:
			return newError(node.GetToken(), "cannot increment %s", val.Type())
		}
	case "--":
		switch v := val.(type) {
		case *object.Integer:
			newVal := &object.Integer{Value: v.Value - 1}
			env.Update(ident.Value, newVal)
			return v
		case *object.Float:
			newVal := &object.Float{Value: v.Value - 1.0}
			env.Update(ident.Value, newVal)
			return v
		default:
			return newError(node.GetToken(), "cannot decrement %s", val.Type())
		}
	default:
		return newError(node.GetToken(), "unknown postfix operator: %s", node.Operator)
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

func evalMinusPrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError(tok, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(tok, operator, left, right)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ):
		return evalFloatInfixExpression(tok, operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(tok, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case left.Type() != right.Type():
		return newError(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
	case "%":
		if rightVal == 0 {
			return newError(tok, "division by zero")
		}
		return &object.Integer{Value: leftVal % rightVal}
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
	case "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	leftVal := getFloatVal(left)
	rightVal := getFloatVal(right)

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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, object.NewEnclosedEnvironment(env))
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, object.NewEnclosedEnvironment(env))
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

		loopEnv := object.NewEnclosedEnvironment(env)
		result = evalBlockStatement(ws.Body, loopEnv)
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

		loopBodyEnv := object.NewEnclosedEnvironment(forEnv)
		result = evalBlockStatement(fs.Body, loopBodyEnv)
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

func applyFunction(tok token.Token, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(tok, "wrong number of arguments for %s. got=%d, want=%d",
				fn.Name, len(args), len(fn.Parameters))
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

func evalStringIndexExpression(tok token.Token, str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError(tok, "index out of bounds: %d", idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
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
		if integer, ok := obj.(*object.Integer); ok {
			return integer.Value != 0
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

func newError(tok token.Token, format string, args ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}
