package evaluator

import (
	"fmt"
	"io"

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
	return EvalWithOutput(node, env, nil)
}

func EvalWithOutput(node ast.Node, env *object.Environment, out io.Writer) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env, out)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env, out)

	case *ast.ExpressionStatement:
		return EvalWithOutput(node.Expression, env, out)

	case *ast.ReturnStatement:
		val := EvalWithOutput(node.ReturnValue, env, out)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		var val object.Object = NULL
		if node.Value != nil {
			val = EvalWithOutput(node.Value, env, out)
			if isError(val) {
				return val
			}
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ShortVarStatement:
		val := EvalWithOutput(node.Value, env, out)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.AssignmentStatement:
		val := EvalWithOutput(node.Value, env, out)
		if isError(val) {
			return val
		}
		if _, ok := env.Update(node.Name.Value, val); !ok {
			return newErrorAt("identifier not found: %s", node.Token, node.Name.Value)
		}
		return NULL

	case *ast.WhileStatement:
		return evalWhileStatement(node, env, out)

	case *ast.ForStatement:
		return evalForStatement(node, env, out)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := EvalWithOutput(node.Right, env, out)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token)

	case *ast.InfixExpression:
		// Handle short-circuit evaluation for && and ||
		if node.Operator == "&&" {
			left := EvalWithOutput(node.Left, env, out)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := EvalWithOutput(node.Right, env, out)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}
		if node.Operator == "||" {
			left := EvalWithOutput(node.Left, env, out)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := EvalWithOutput(node.Right, env, out)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := EvalWithOutput(node.Left, env, out)
		if isError(left) {
			return left
		}
		right := EvalWithOutput(node.Right, env, out)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env, out)

	case *ast.Identifier:
		return evalIdentifier(node, env, out)

	case *ast.FunctionLiteral:
		fn := &object.Function{
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
		function := EvalWithOutput(node.Function, env, out)
		if isError(function) {
			return function
		}
		args := evalExpressionsWithOutput(node.Arguments, env, out)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunctionWithOutput(function, args, node.Token, out)

	case *ast.ArrayLiteral:
		elems := evalExpressionsWithOutput(node.Elements, env, out)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.IndexExpression:
		left := EvalWithOutput(node.Left, env, out)
		if isError(left) {
			return left
		}
		index := EvalWithOutput(node.Index, env, out)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, node.Token)

	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment, out io.Writer) object.Object {
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = EvalWithOutput(statement, env, out)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	// Check if a zero-arg main function is defined in statements and auto-invoke
	if hasMainFunction(program) {
		if mainObj, ok := env.Get("main"); ok {
			if mainFn, ok := mainObj.(*object.Function); ok {
				return applyFunctionWithOutput(mainFn, []object.Object{}, token.Token{}, out)
			}
		}
	}

	return result
}

func hasMainFunction(program *ast.Program) bool {
	for _, stmt := range program.Statements {
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if fnLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				if fnLit.Name != nil && fnLit.Name.Value == "main" && len(fnLit.Parameters) == 0 {
					return true
				}
			}
		}
	}
	return false
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, out io.Writer) object.Object {
	var result object.Object = NULL

	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = EvalWithOutput(statement, blockEnv, out)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment, out io.Writer) object.Object {
	var result object.Object = NULL

	for {
		condition := EvalWithOutput(ws.Condition, env, out)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = evalBlockStatement(ws.Body, env, out)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment, out io.Writer) object.Object {
	var result object.Object = NULL
	forEnv := object.NewEnclosedEnvironment(env)

	if fs.Init != nil {
		initRes := EvalWithOutput(fs.Init, forEnv, out)
		if isError(initRes) {
			return initRes
		}
	}

	for {
		if fs.Condition != nil {
			condition := EvalWithOutput(fs.Condition, forEnv, out)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = evalBlockStatement(fs.Body, forEnv, out)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}

		if fs.Post != nil {
			postRes := EvalWithOutput(fs.Post, forEnv, out)
			if isError(postRes) {
				return postRes
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

func evalPrefixExpression(operator string, right object.Object, tok token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, tok)
	default:
		return newErrorAt("unknown operator: %s%s", tok, operator, right.Type())
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
		if isTruthy(right) {
			return FALSE
		}
		return TRUE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		val := right.(*object.Integer).Value
		return &object.Integer{Value: -val}
	case object.FLOAT_OBJ:
		val := right.(*object.Float).Value
		return &object.Float{Value: -val}
	default:
		return newErrorAt("unknown operator: -%s", tok, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ):
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return newErrorAt("type mismatch: %s %s %s", tok, left.Type(), operator, right.Type())
	default:
		return newErrorAt("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
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
			return newErrorAt("division by zero", tok)
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
		return newErrorAt("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
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
			return newErrorAt("division by zero", tok)
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
		return newErrorAt("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
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
		return newErrorAt("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newErrorAt("postfix operator %s requires identifier", node.Token, node.Operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newErrorAt("identifier not found: %s", node.Token, ident.Value)
	}

	switch val := val.(type) {
	case *object.Integer:
		if node.Operator == "++" {
			env.Update(ident.Value, &object.Integer{Value: val.Value + 1})
			return val
		} else if node.Operator == "--" {
			env.Update(ident.Value, &object.Integer{Value: val.Value - 1})
			return val
		}
	case *object.Float:
		if node.Operator == "++" {
			env.Update(ident.Value, &object.Float{Value: val.Value + 1})
			return val
		} else if node.Operator == "--" {
			env.Update(ident.Value, &object.Float{Value: val.Value - 1})
			return val
		}
	}

	return newErrorAt("postfix operator %s not supported for type %s", node.Token, node.Operator, val.Type())
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, out io.Writer) object.Object {
	condition := EvalWithOutput(ie.Condition, env, out)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return evalBlockStatement(ie.Consequence, env, out)
	} else if ie.Alternative != nil {
		return evalBlockStatement(ie.Alternative, env, out)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment, out io.Writer) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if out != nil {
		bMap := RegisterBuiltinsWithWriter(out)
		if builtin, ok := bMap[node.Value]; ok {
			return builtin
		}
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorAt("identifier not found: %s", node.Token, node.Value)
}

func evalExpressionsWithOutput(exps []ast.Expression, env *object.Environment, out io.Writer) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := EvalWithOutput(e, env, out)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, tok)
	default:
		return newErrorAt("index operator not supported: %s[%s]", tok, left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, tok token.Token) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newErrorAt("index out of bounds: index %d, len %d", tok, idx, len(arrayObject.Elements))
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, tok token.Token) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newErrorAt("index out of bounds: index %d, len %d", tok, idx, len(strObject.Value))
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func applyFunctionWithOutput(fn object.Object, args []object.Object, tok token.Token, out io.Writer) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := EvalWithOutput(fn.Body, extendedEnv, out)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newErrorAt("not a function: %s", tok, fn.Type())
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

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}
	switch left := left.(type) {
	case *object.Integer:
		return left.Value == right.(*object.Integer).Value
	case *object.Float:
		return left.Value == right.(*object.Float).Value
	case *object.Boolean:
		return left.Value == right.(*object.Boolean).Value
	case *object.String:
		return left.Value == right.(*object.String).Value
	case *object.Null:
		return true
	default:
		return left == right
	}
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
		if flt, ok := obj.(*object.Float); ok {
			return flt.Value != 0
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

func newErrorAt(format string, tok token.Token, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}
