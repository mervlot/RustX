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

func Eval(node ast.Node, env *object.Environment) object.Object {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return NULL

	case *ast.WhileStatement:
		return evalWhileStatement(n, env)

	case *ast.ForStatement:
		return evalForStatement(n, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(n.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(n, left, index)

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n, n.Operator, right)

	case *ast.InfixExpression:
		if n.Operator == "&&" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}
		if n.Operator == "||" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n, n.Operator, left, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(n, env)

	case *ast.AssignmentExpression:
		return evalAssignmentExpression(n, env)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if n.Name != nil {
			env.Set(n.Name.Value, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(n, function, args)
	}

	return nil
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(node ast.Node, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
	default:
		return newError("unknown operator: %s%s", node.GetToken().Line, node.GetToken().Column, operator, right.Type())
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

func evalMinusPrefixOperatorExpression(node ast.Node, right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		val := right.(*object.Integer).Value
		return &object.Integer{Value: -val}
	case object.FLOAT_OBJ:
		val := right.(*object.Float).Value
		return &object.Float{Value: -val}
	default:
		return newError("unknown operator: -%s", node.GetToken().Line, node.GetToken().Column, right.Type())
	}
}

func evalInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, operator, left.(*object.Float), right.(*object.Float))
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		promotedRight := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, left.(*object.Float), promotedRight)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		promotedLeft := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, promotedLeft, right.(*object.Float))
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, operator, left.(*object.String), right.(*object.String))
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newError("unknown operator: %s %s %s", node.GetToken().Line, node.GetToken().Column, left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(node ast.Node, operator string, left, right *object.Integer) object.Object {
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
			return newError("division by zero", node.GetToken().Line, node.GetToken().Column)
		}
		return &object.Integer{Value: leftVal / rightVal}
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
		return newError("unknown operator: %s %s %s", node.GetToken().Line, node.GetToken().Column, left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(node ast.Node, operator string, left, right *object.Float) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero", node.GetToken().Line, node.GetToken().Column)
		}
		return &object.Float{Value: leftVal / rightVal}
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
		return newError("unknown operator: %s %s %s", node.GetToken().Line, node.GetToken().Column, left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(node ast.Node, operator string, left, right *object.String) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", node.GetToken().Line, node.GetToken().Column, left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		if (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ) ||
			(left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) {
			var lVal, rVal float64
			if left.Type() == object.INTEGER_OBJ {
				lVal = float64(left.(*object.Integer).Value)
			} else {
				lVal = left.(*object.Float).Value
			}
			if right.Type() == object.INTEGER_OBJ {
				rVal = float64(right.(*object.Integer).Value)
			} else {
				rVal = right.(*object.Float).Value
			}
			return lVal == rVal
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
		result = Eval(ws.Body, env)
		if isError(result) {
			return result
		}
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv
		}
	}
	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	loopEnv := object.NewEnclosedEnvironment(env)

	if fs.Init != nil {
		initVal := Eval(fs.Init, loopEnv)
		if isError(initVal) {
			return initVal
		}
	}

	var result object.Object = NULL

	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, loopEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fs.Body, loopEnv)
		if isError(result) {
			return result
		}
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv
		}

		if fs.Post != nil {
			postVal := Eval(fs.Post, loopEnv)
			if isError(postVal) {
				return postVal
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

	return newError("identifier not found: %s", node.Token.Line, node.Token.Column, node.Value)
}

func evalPostfixExpression(pe *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := pe.Left.(*ast.Identifier)
	if !ok {
		return newError("postfix operators can only be applied to identifiers", pe.Token.Line, pe.Token.Column)
	}

	val, found := env.Get(ident.Value)
	if !found {
		return newError("identifier not found: %s", pe.Token.Line, pe.Token.Column, ident.Value)
	}

	switch val := val.(type) {
	case *object.Integer:
		var newVal int64
		if pe.Operator == "++" {
			newVal = val.Value + 1
		} else {
			newVal = val.Value - 1
		}
		ret := &object.Integer{Value: newVal}
		env.Update(ident.Value, ret)
		return val
	case *object.Float:
		var newVal float64
		if pe.Operator == "++" {
			newVal = val.Value + 1.0
		} else {
			newVal = val.Value - 1.0
		}
		ret := &object.Float{Value: newVal}
		env.Update(ident.Value, ret)
		return val
	default:
		return newError("postfix operators can only be applied to numbers, got %s", pe.Token.Line, pe.Token.Column, val.Type())
	}
}

func evalAssignmentExpression(ae *ast.AssignmentExpression, env *object.Environment) object.Object {
	val := Eval(ae.Value, env)
	if isError(val) {
		return val
	}

	if ident, ok := ae.Left.(*ast.Identifier); ok {
		if ae.Operator == ":=" {
			env.Set(ident.Value, val)
			return val
		} else {
			_, updated := env.Update(ident.Value, val)
			if !updated {
				return newError("identifier not found: %s", ae.Token.Line, ae.Token.Column, ident.Value)
			}
			return val
		}
	}

	if ie, ok := ae.Left.(*ast.IndexExpression); ok {
		if ae.Operator == ":=" {
			return newError("cannot use shorthand assignment := on index expression", ae.Token.Line, ae.Token.Column)
		}
		leftVal := Eval(ie.Left, env)
		if isError(leftVal) {
			return leftVal
		}
		indexVal := Eval(ie.Index, env)
		if isError(indexVal) {
			return indexVal
		}

		arr, isArr := leftVal.(*object.Array)
		if !isArr {
			return newError("cannot assign index of non-array: %s", ae.Token.Line, ae.Token.Column, leftVal.Type())
		}
		idx, isInt := indexVal.(*object.Integer)
		if !isInt {
			return newError("index must be integer, got %s", ae.Token.Line, ae.Token.Column, indexVal.Type())
		}

		if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
			return newError("index out of bounds: %d (array size %d)", ae.Token.Line, ae.Token.Column, idx.Value, len(arr.Elements))
		}

		arr.Elements[idx.Value] = val
		return val
	}

	return newError("invalid assignment target", ae.Token.Line, ae.Token.Column)
}

func evalIndexExpression(ie *ast.IndexExpression, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arr := left.(*object.Array)
		idx := index.(*object.Integer).Value
		if idx < 0 || idx >= int64(len(arr.Elements)) {
			return newError("index out of bounds: %d (array size %d)", ie.Token.Line, ie.Token.Column, idx, len(arr.Elements))
		}
		return arr.Elements[idx]
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		str := left.(*object.String).Value
		idx := index.(*object.Integer).Value
		if idx < 0 || idx >= int64(len(str)) {
			return newError("index out of bounds: %d (string size %d)", ie.Token.Line, ie.Token.Column, idx, len(str))
		}
		return &object.String{Value: string(str[idx])}
	default:
		return newError("index operator not supported for %s[%s]", ie.Token.Line, ie.Token.Column, left.Type(), index.Type())
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

func applyFunction(node ast.Node, fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		if len(args) != len(f.Parameters) {
			return newError("wrong number of arguments. got=%d, want=%d", node.GetToken().Line, node.GetToken().Column, len(args), len(f.Parameters))
		}
		extendedEnv := object.NewEnclosedEnvironment(f.Env)
		for paramIdx, param := range f.Parameters {
			extendedEnv.Set(param.Name.Value, args[paramIdx])
		}
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(node.GetToken().Line, node.GetToken().Column, args...)

	default:
		return newError("not a function: %s", node.GetToken().Line, node.GetToken().Column, fn.Type())
	}
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
		switch val := obj.(type) {
		case *object.Boolean:
			return val.Value
		case *object.Integer:
			return val.Value != 0
		case *object.Float:
			return val.Value != 0.0
		default:
			return true
		}
	}
}

func newError(format string, line, col int, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...), Line: line, Column: col}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
