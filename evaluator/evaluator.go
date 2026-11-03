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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message:  fmt.Sprintf(format, a...),
		Position: token.Position{},
	}
}

func newPosError(pos token.Position, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message:  fmt.Sprintf(format, a...),
		Position: pos,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// Eval evaluates an AST node within the given environment.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
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

	case *ast.ShortVarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.AssignStatement:
		return evalAssignStatement(node, env)

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
		return evalIndexExpression(left, index, node.Token.Position)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Position)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		if node.Operator == "&&" {
			return evalLogicalAnd(node.Left, node.Right, env)
		}
		if node.Operator == "||" {
			return evalLogicalOr(node.Left, node.Right, env)
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token.Position)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Name:       node.Name,
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		// If named function, register in current environment
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
		return applyFunction(function, args, node.Token.Position)
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

func evalAssignStatement(node *ast.AssignStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	switch left := node.Left.(type) {
	case *ast.Identifier:
		if !env.Update(left.Value, val) {
			// If not in any outer scope, set in current scope
			env.Set(left.Value, val)
		}
		return val
	case *ast.IndexExpression:
		arrObj := Eval(left.Left, env)
		if isError(arrObj) {
			return arrObj
		}
		idxObj := Eval(left.Index, env)
		if isError(idxObj) {
			return idxObj
		}
		arr, ok := arrObj.(*object.Array)
		if !ok {
			return newPosError(left.Token.Position, "index assignment only supported on arrays, got %s", arrObj.Type())
		}
		idx, ok := idxObj.(*object.Integer)
		if !ok {
			return newPosError(left.Token.Position, "array index must be integer, got %s", idxObj.Type())
		}
		if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
			return newPosError(left.Token.Position, "index out of bounds: len=%d, index=%d", len(arr.Elements), idx.Value)
		}
		arr.Elements[idx.Value] = val
		return val
	default:
		return newPosError(node.Token.Position, "invalid assignment target: %s", node.Left.String())
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newPosError(node.Token.Position, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object, pos token.Position) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, pos)
	default:
		return newPosError(pos, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, pos token.Position) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newPosError(pos, "unknown operator: -%s", right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newPosError(node.Token.Position, "operator %s requires identifier", node.Operator)
	}

	val, exists := env.Get(ident.Value)
	if !exists {
		return newPosError(node.Token.Position, "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		oldVal := v.Value
		newVal := oldVal
		if node.Operator == "++" {
			newVal++
		} else if node.Operator == "--" {
			newVal--
		}
		updatedObj := &object.Integer{Value: newVal}
		if !env.Update(ident.Value, updatedObj) {
			env.Set(ident.Value, updatedObj)
		}
		return &object.Integer{Value: oldVal}
	case *object.Float:
		oldVal := v.Value
		newVal := oldVal
		if node.Operator == "++" {
			newVal += 1.0
		} else if node.Operator == "--" {
			newVal -= 1.0
		}
		updatedObj := &object.Float{Value: newVal}
		if !env.Update(ident.Value, updatedObj) {
			env.Set(ident.Value, updatedObj)
		}
		return &object.Float{Value: oldVal}
	default:
		return newPosError(node.Token.Position, "operator %s not supported on type %s", node.Operator, val.Type())
	}
}

func evalLogicalAnd(leftExp, rightExp ast.Expression, env *object.Environment) object.Object {
	left := Eval(leftExp, env)
	if isError(left) {
		return left
	}
	if !isTruthy(left) {
		return FALSE
	}
	right := Eval(rightExp, env)
	if isError(right) {
		return right
	}
	return nativeBoolToBooleanObject(isTruthy(right))
}

func evalLogicalOr(leftExp, rightExp ast.Expression, env *object.Environment) object.Object {
	left := Eval(leftExp, env)
	if isError(left) {
		return left
	}
	if isTruthy(left) {
		return TRUE
	}
	right := Eval(rightExp, env)
	if isError(right) {
		return right
	}
	return nativeBoolToBooleanObject(isTruthy(right))
}

func evalInfixExpression(operator string, left, right object.Object, pos token.Position) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, pos)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ) &&
		(left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ):
		return evalFloatInfixExpression(operator, left, right, pos)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, pos)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return newPosError(pos, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newPosError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left == right {
		return true
	}
	switch l := left.(type) {
	case *object.Integer:
		if r, ok := right.(*object.Integer); ok {
			return l.Value == r.Value
		}
	case *object.Float:
		if r, ok := right.(*object.Float); ok {
			return l.Value == r.Value
		}
	case *object.Boolean:
		if r, ok := right.(*object.Boolean); ok {
			return l.Value == r.Value
		}
	case *object.String:
		if r, ok := right.(*object.String); ok {
			return l.Value == r.Value
		}
	}
	return false
}

func evalIntegerInfixExpression(operator string, left, right object.Object, pos token.Position) object.Object {
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
			return newPosError(pos, "division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newPosError(pos, "modulo by zero")
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
	default:
		return newPosError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, pos token.Position) object.Object {
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
			return newPosError(pos, "division by zero")
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
		return newPosError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, pos token.Position) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	default:
		return newPosError(pos, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	forEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		initRes := Eval(fe.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fe.Condition != nil {
			condition := Eval(fe.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, forEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}

		if fe.Post != nil {
			postRes := Eval(fe.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return result
}

func evalIndexExpression(left, index object.Object, pos token.Position) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, pos)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, pos)
	default:
		return newPosError(pos, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, pos token.Position) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newPosError(pos, "array index out of bounds: len=%d, index=%d", len(arrayObject.Elements), idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, pos token.Position) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newPosError(pos, "string index out of bounds: len=%d, index=%d", len(strObject.Value), idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func applyFunction(fn object.Object, args []object.Object, pos token.Position) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		if len(args) != len(function.Parameters) {
			return newPosError(pos, "wrong number of arguments: expected %d, got %d", len(function.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return function.Fn(args...)

	default:
		return newPosError(pos, "not a function: %s", fn.Type())
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
			return floatObj.Value != 0.0
		}
		if strObj, ok := obj.(*object.String); ok {
			return strObj.Value != ""
		}
		return true
	}
}
