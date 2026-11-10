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
	case *ast.ShortVarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL
	case *ast.AssignStatement:
		return evalAssignStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)
	case *ast.InfixExpression:
		// Short-circuit logical operators
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
		return evalInfixExpression(node.Operator, left, right, node.Token)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
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
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, node.Token)
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
		return evalIndexExpression(left, index, node.Token)
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

	switch target := node.Left.(type) {
	case *ast.Identifier:
		updated := env.Update(target.Value, val)
		if !updated {
			return newError(node.Token, "undefined variable %q", target.Value)
		}
		return NULL
	case *ast.IndexExpression:
		left := Eval(target.Left, env)
		if isError(left) {
			return left
		}
		idx := Eval(target.Index, env)
		if isError(idx) {
			return idx
		}
		arr, ok := left.(*object.Array)
		if !ok {
			return newError(target.Token, "index assignment only supported on arrays, got %s", left.Type())
		}
		idxInt, ok := idx.(*object.Integer)
		if !ok {
			return newError(target.Token, "index must be integer, got %s", idx.Type())
		}
		if idxInt.Value < 0 || int(idxInt.Value) >= len(arr.Elements) {
			return newError(target.Token, "index out of bounds: index %d, length %d", idxInt.Value, len(arr.Elements))
		}
		arr.Elements[idxInt.Value] = val
		return NULL
	default:
		return newError(node.Token, "invalid assignment target")
	}
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
		if right.Type() == object.INTEGER_OBJ && right.(*object.Integer).Value == 0 {
			return TRUE
		}
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError(tok, "unknown operator: -%s", right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.Token, "operator %s only valid on variables", node.Operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(node.Token, "undefined variable %q", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		oldVal := &object.Integer{Value: v.Value}
		if node.Operator == "++" {
			env.Update(ident.Value, &object.Integer{Value: v.Value + 1})
		} else if node.Operator == "--" {
			env.Update(ident.Value, &object.Integer{Value: v.Value - 1})
		}
		return oldVal
	case *object.Float:
		oldVal := &object.Float{Value: v.Value}
		if node.Operator == "++" {
			env.Update(ident.Value, &object.Float{Value: v.Value + 1.0})
		} else if node.Operator == "--" {
			env.Update(ident.Value, &object.Float{Value: v.Value - 1.0})
		}
		return oldVal
	default:
		return newError(node.Token, "operator %s not supported for type %s", node.Operator, val.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ) &&
		(left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ):
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
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
			return newError(tok, "division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newError(tok, "modulo by zero")
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
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
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

	// Evaluate Init if present
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.Token, "identifier not found: %s", node.Value)
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

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, tok)
	default:
		return newError(tok, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, tok token.Token) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError(tok, "index out of bounds: index %d, length %d", idx, len(arrayObject.Elements))
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, tok token.Token) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError(tok, "index out of bounds: index %d, length %d", idx, len(stringObject.Value))
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func applyFunction(fn object.Object, args []object.Object, tok token.Token) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(args) != len(fn.Parameters) {
			return newError(tok, "wrong number of arguments: expected %d, got %d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		res := fn.Fn(args...)
		if err, ok := res.(*object.Error); ok && err.Line == 0 {
			err.Line = tok.Line
			err.Column = tok.Column
		}
		return res
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
		if obj.Type() == object.INTEGER_OBJ && obj.(*object.Integer).Value == 0 {
			return false
		}
		return true
	}
}

func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
