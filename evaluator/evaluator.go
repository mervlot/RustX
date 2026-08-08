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
		return NULL
	}

	switch n := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.VarStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return NULL

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

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
		tok := n.GetToken()
		return evalIndexExpression(left, index, tok.Line, tok.Column)

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		tok := n.GetToken()
		return evalPrefixExpression(n.Operator, right, tok.Line, tok.Column)

	case *ast.InfixExpression:
		tok := n.GetToken()
		// Logical AND and OR with short-circuiting
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
		return evalInfixExpression(n.Operator, left, right, tok.Line, tok.Column)

	case *ast.PostfixExpression:
		tok := n.GetToken()
		return evalPostfixExpression(n, env, tok.Line, tok.Column)

	case *ast.AssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		ok := env.Update(n.Name.Value, val)
		if !ok {
			tok := n.GetToken()
			return newError(tok.Line, tok.Column, "variable %s is not defined", n.Name.Value)
		}
		return val

	case *ast.ShortAssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return val

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		funcObj := &object.Function{
			Name:       n.Name,
			Parameters: n.Parameters,
			Body:       n.Body,
			Env:        env,
		}
		// Register in env if it's a named function
		if n.Name != nil {
			env.Set(n.Name.Value, funcObj)
		}
		return funcObj

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		tok := n.GetToken()
		return applyFunction(function, args, tok.Line, tok.Column)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // initialized to NULL to ensure consistent behavior on empty files

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, stmt := range block.Statements {
		result = Eval(stmt, blockEnv)

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

func evalPrefixExpression(operator string, right object.Object, line, col int) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, line, col)
	default:
		return newError(line, col, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, line, col int) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError(line, col, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, line, col int) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer), line, col)
	}
	if left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ {
		return evalFloatInfixExpression(operator, left.(*object.Float), right.(*object.Float), line, col)
	}
	// Mixed arithmetic (Integer & Float)
	if (left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ) ||
		(left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ) {
		return evalMixedInfixExpression(operator, left, right, line, col)
	}

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String), line, col)
	}

	if operator == "==" {
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	}
	if operator == "!=" {
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	}

	return newError(line, col, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
}

func evalIntegerInfixExpression(op string, left, right *object.Integer, line, col int) object.Object {
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
			return newError(line, col, "division by zero")
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
	}
	return newError(line, col, "unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalFloatInfixExpression(op string, left, right *object.Float, line, col int) object.Object {
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
			return newError(line, col, "division by zero")
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
	}
	return newError(line, col, "unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalMixedInfixExpression(op string, left, right object.Object, line, col int) object.Object {
	var lFloat, rFloat float64

	if left.Type() == object.INTEGER_OBJ {
		lFloat = float64(left.(*object.Integer).Value)
	} else {
		lFloat = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rFloat = float64(right.(*object.Integer).Value)
	} else {
		rFloat = right.(*object.Float).Value
	}

	return evalFloatInfixExpression(op, &object.Float{Value: lFloat}, &object.Float{Value: rFloat}, line, col)
}

func evalStringInfixExpression(op string, left, right *object.String, line, col int) object.Object {
	if op == "+" {
		return &object.String{Value: left.Value + right.Value}
	}
	if op == "==" {
		return nativeBoolToBooleanObject(left.Value == right.Value)
	}
	if op == "!=" {
		return nativeBoolToBooleanObject(left.Value != right.Value)
	}
	return newError(line, col, "unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalObjectEquality(left object.Object, right object.Object) bool {
	if left.Type() != right.Type() {
		// Mixed numeric equality
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

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment, line, col int) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(line, col, "postfix operator can only be applied to identifier")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(line, col, "identifier not found: %s", ident.Value)
	}

	switch node.Operator {
	case "++":
		switch v := val.(type) {
		case *object.Integer:
			newObj := &object.Integer{Value: v.Value + 1}
			env.Update(ident.Value, newObj)
			return v // postfix returns original value
		case *object.Float:
			newObj := &object.Float{Value: v.Value + 1.0}
			env.Update(ident.Value, newObj)
			return v
		default:
			return newError(line, col, "cannot increment non-numeric type %s", val.Type())
		}
	case "--":
		switch v := val.(type) {
		case *object.Integer:
			newObj := &object.Integer{Value: v.Value - 1}
			env.Update(ident.Value, newObj)
			return v
		case *object.Float:
			newObj := &object.Float{Value: v.Value - 1.0}
			env.Update(ident.Value, newObj)
			return v
		default:
			return newError(line, col, "cannot decrement non-numeric type %s", val.Type())
		}
	}

	return newError(line, col, "unknown postfix operator: %s", node.Operator)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	tok := node.GetToken()
	return newError(tok.Line, tok.Column, "identifier not found: %s", node.Value)
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

func evalIndexExpression(left, index object.Object, line, col int) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left.(*object.Array), index.(*object.Integer), line, col)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left.(*object.String), index.(*object.Integer), line, col)
	default:
		return newError(line, col, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array *object.Array, index *object.Integer, line, col int) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return newError(line, col, "array index out of bounds: %d (len %d)", idx, len(array.Elements))
	}

	return array.Elements[idx]
}

func evalStringIndexExpression(str *object.String, index *object.Integer, line, col int) object.Object {
	idx := index.Value
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return newError(line, col, "string index out of bounds: %d (len %d)", idx, len(str.Value))
	}

	return &object.String{Value: string(str.Value[idx])}
}

func applyFunction(fn object.Object, args []object.Object, line, col int) object.Object {
	switch f := fn.(type) {
	case *object.Builtin:
		return f.Fn(line, col, args...)

	case *object.Function:
		if len(args) != len(f.Parameters) {
			return newError(line, col, "wrong number of arguments. got=%d, want=%d", len(args), len(f.Parameters))
		}
		extendedEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	default:
		return newError(line, col, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	// Register named function inside itself for support of recursion
	if fn.Name != nil {
		env.Set(fn.Name.Value, fn)
	}

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

func newError(line, col int, format string, args ...interface{}) *object.Error {
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
