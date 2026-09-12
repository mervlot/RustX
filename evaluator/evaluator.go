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
		if val == nil {
			val = NULL
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.ShortVarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return val

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		_, ok := env.Update(node.Name.Value, val)
		if !ok {
			tok := node.GetToken()
			return newErrorWithPos(tok.Line, tok.Column, "identifier not found: %s", node.Name.Value)
		}
		return val

	case *ast.ReturnStatement:
		var val object.Object = NULL
		if node.ReturnValue != nil {
			val = Eval(node.ReturnValue, env)
			if isError(val) {
				return val
			}
		}
		return &object.ReturnValue{Value: val}

	case *ast.FunctionStatement:
		fn := &object.Function{
			Parameters: node.Parameters,
			ReturnType: node.ReturnType,
			Body:       node.Body,
			Env:        env,
		}
		if node.Name != nil {
			fn.Name = node.Name.Value
			env.Set(node.Name.Value, fn)
		}
		return fn

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

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node)

	case *ast.InfixExpression:
		// Handle short-circuiting for logical operators
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
		return evalInfixExpression(node.Operator, left, right, node)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

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

	// Automatic entrypoint invocation if main() function is defined
	if mainObj, ok := env.Get("main"); ok {
		if mainFn, ok := mainObj.(*object.Function); ok {
			if len(mainFn.Parameters) == 0 {
				mainRes := applyFunction(mainFn, []object.Object{}, program)
				if isError(mainRes) {
					return mainRes
				}
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalWhileStatement(stmt *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		cond := Eval(stmt.Condition, env)
		if isError(cond) {
			return cond
		}
		if !isTruthy(cond) {
			break
		}

		result = evalBlockStatement(stmt.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(stmt *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)
	var result object.Object = NULL

	if stmt.Init != nil {
		initRes := Eval(stmt.Init, loopEnv)
		if isError(initRes) {
			return initRes
		}
	}

	for {
		if stmt.Condition != nil {
			cond := Eval(stmt.Condition, loopEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
		}

		if stmt.Body != nil {
			result = evalBlockStatement(stmt.Body, loopEnv)
			if result != nil {
				rt := result.Type()
				if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
					return result
				}
			}
		}

		if stmt.Post != nil {
			postRes := Eval(stmt.Post, loopEnv)
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

	tok := node.GetToken()
	return newErrorWithPos(tok.Line, tok.Column, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	tok := node.GetToken()
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, tok)
	default:
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: %s%s", operator, right.Type())
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
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	tok := node.GetToken()
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return newErrorWithPos(tok.Line, tok.Column, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
			return newErrorWithPos(tok.Line, tok.Column, "division by zero")
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
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	var leftVal, rightVal float64

	if l, ok := left.(*object.Integer); ok {
		leftVal = float64(l.Value)
	} else if l, ok := left.(*object.Float); ok {
		leftVal = l.Value
	} else {
		return newErrorWithPos(tok.Line, tok.Column, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	if r, ok := right.(*object.Integer); ok {
		rightVal = float64(r.Value)
	} else if r, ok := right.(*object.Float); ok {
		rightVal = r.Value
	} else {
		return newErrorWithPos(tok.Line, tok.Column, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
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
			return newErrorWithPos(tok.Line, tok.Column, "division by zero")
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
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newErrorWithPos(tok.Line, tok.Column, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		tok := node.GetToken()
		return newErrorWithPos(tok.Line, tok.Column, "postfix operand must be an identifier")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		tok := node.GetToken()
		return newErrorWithPos(tok.Line, tok.Column, "identifier not found: %s", ident.Value)
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
		return v
	case *object.Float:
		var newVal float64
		if node.Operator == "++" {
			newVal = v.Value + 1.0
		} else {
			newVal = v.Value - 1.0
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return v
	default:
		tok := node.GetToken()
		return newErrorWithPos(tok.Line, tok.Column, "postfix operator %s not supported for type %s", node.Operator, val.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return evalBlockStatement(ie.Consequence, env)
	} else if ie.Alternative != nil {
		switch alt := ie.Alternative.(type) {
		case *ast.BlockStatement:
			return evalBlockStatement(alt, env)
		default:
			return Eval(alt, env)
		}
	} else {
		return NULL
	}
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

func evalIndexExpression(left, index object.Object, node ast.Node) object.Object {
	tok := node.GetToken()
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, tok)
	default:
		return newErrorWithPos(tok.Line, tok.Column, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, tok token.Token) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newErrorWithPos(tok.Line, tok.Column, "array index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, tok token.Token) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newErrorWithPos(tok.Line, tok.Column, "string index out of bounds: %d", idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func applyFunction(fn object.Object, args []object.Object, node ast.Node) object.Object {
	tok := node.GetToken()
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newErrorWithPos(tok.Line, tok.Column, "wrong number of arguments for function %s: expected %d, got %d", fn.Name, len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newErrorWithPos(tok.Line, tok.Column, "not a function: %s", fn.Type())
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
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
		if floatVal, ok := obj.(*object.Float); ok {
			return floatVal.Value != 0
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

func newErrorWithPos(line, col int, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    line,
		Column:  col,
	}
}
