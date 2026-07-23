package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/token"
	"github.com/username/rustx/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
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
		var val object.Object = NULL
		if n.Value != nil {
			val = Eval(n.Value, env)
			if isError(val) {
				return val
			}
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
	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.ArrayLiteral:
		elems := evalExpressions(n.Elements, env)
		if len(elems) > 0 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.IndexExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		idx := Eval(n.Index, env)
		if isError(idx) {
			return idx
		}
		return evalIndexExpression(left, idx, n.GetToken())

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right, n.GetToken())

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
		return evalInfixExpression(n.Operator, left, right, n.GetToken())

	case *ast.AssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}

		// 1. Variable Assignment
		ident, ok := n.Left.(*ast.Identifier)
		if ok {
			updated := env.Update(ident.Value, val)
			if !updated {
				return newError(ident.GetToken(), "identifier %s not defined", ident.Value)
			}
			return val
		}

		// 2. Index Assignment
		idxExp, ok := n.Left.(*ast.IndexExpression)
		if ok {
			leftVal := Eval(idxExp.Left, env)
			if isError(leftVal) {
				return leftVal
			}
			idxVal := Eval(idxExp.Index, env)
			if isError(idxVal) {
				return idxVal
			}
			arr, ok := leftVal.(*object.Array)
			if !ok {
				return newError(idxExp.Left.GetToken(), "index assignment target must be an array, got %s", leftVal.Type())
			}
			idxInt, ok := idxVal.(*object.Integer)
			if !ok {
				return newError(idxExp.Index.GetToken(), "array index must be an integer, got %s", idxVal.Type())
			}
			idx := idxInt.Value
			if idx < 0 || idx >= int64(len(arr.Elements)) {
				return newError(idxExp.Index.GetToken(), "index out of bounds: %d, length is %d", idx, len(arr.Elements))
			}
			arr.Elements[idx] = val
			return val
		}

		return newError(n.Left.GetToken(), "invalid assignment target")

	case *ast.ColonAssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return newError(n.Left.GetToken(), "left side of shorthand assignment must be an identifier")
		}
		env.Set(ident.Value, val)
		return val

	case *ast.PostfixExpression:
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return newError(n.Left.GetToken(), "postfix operator %s must be applied to an identifier", n.Operator)
		}
		val, ok := env.Get(ident.Value)
		if !ok {
			return newError(n.Left.GetToken(), "identifier %s not defined", ident.Value)
		}

		if n.Operator == "++" {
			switch v := val.(type) {
			case *object.Integer:
				newVal := &object.Integer{Value: v.Value + 1}
				env.Update(ident.Value, newVal)
				return newVal
			case *object.Float:
				newVal := &object.Float{Value: v.Value + 1.0}
				env.Update(ident.Value, newVal)
				return newVal
			default:
				return newError(n.Left.GetToken(), "postfix ++ can only be applied to numeric type, got %s", val.Type())
			}
		} else if n.Operator == "--" {
			switch v := val.(type) {
			case *object.Integer:
				newVal := &object.Integer{Value: v.Value - 1}
				env.Update(ident.Value, newVal)
				return newVal
			case *object.Float:
				newVal := &object.Float{Value: v.Value - 1.0}
				env.Update(ident.Value, newVal)
				return newVal
			default:
				return newError(n.Left.GetToken(), "postfix -- can only be applied to numeric type, got %s", val.Type())
			}
		}
		return newError(n.GetToken(), "unknown postfix operator: %s", n.Operator)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Parameters: n.Parameters,
			Body:       n.Body,
			Env:        env,
		}
		if n.Name != "" {
			env.Set(n.Name, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) > 0 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, n.GetToken())
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = Eval(statement, env)

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

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch l := left.(type) {
	case *object.Array:
		idx, ok := index.(*object.Integer)
		if !ok {
			return newError(tok, "array index must be an integer, got %s", index.Type())
		}
		if idx.Value < 0 || idx.Value >= int64(len(l.Elements)) {
			return newError(tok, "index out of bounds: %d, length is %d", idx.Value, len(l.Elements))
		}
		return l.Elements[idx.Value]
	case *object.String:
		idx, ok := index.(*object.Integer)
		if !ok {
			return newError(tok, "string index must be an integer, got %s", index.Type())
		}
		if idx.Value < 0 || idx.Value >= int64(len(l.Value)) {
			return newError(tok, "index out of bounds: %d, length is %d", idx.Value, len(l.Value))
		}
		return &object.String{Value: string(l.Value[idx.Value])}
	default:
		return newError(tok, "index operator not supported on %s", left.Type())
	}
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

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	// Division by Zero check
	if operator == "/" {
		switch r := right.(type) {
		case *object.Integer:
			if r.Value == 0 {
				return newError(tok, "division by zero")
			}
		case *object.Float:
			if r.Value == 0.0 {
				return newError(tok, "division by zero")
			}
		}
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Promote left to float
		lFloat := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, lFloat, right, tok)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Promote right to float
		rFloat := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left, rFloat, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
	case operator == "==":
		return evalObjectEquality(left, right, "==", tok)
	case operator == "!=":
		return evalObjectEquality(left, right, "!=", tok)
	default:
		return newError(tok, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
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

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
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
	default:
		return newError(tok, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object, op string, tok token.Token) object.Object {
	var eq bool

	switch l := left.(type) {
	case *object.Integer:
		switch r := right.(type) {
		case *object.Integer:
			eq = (l.Value == r.Value)
		case *object.Float:
			eq = (float64(l.Value) == r.Value)
		default:
			return newError(tok, "type mismatch: %s %s %s", left.Type(), op, right.Type())
		}
	case *object.Float:
		switch r := right.(type) {
		case *object.Integer:
			eq = (l.Value == float64(r.Value))
		case *object.Float:
			eq = (l.Value == r.Value)
		default:
			return newError(tok, "type mismatch: %s %s %s", left.Type(), op, right.Type())
		}
	case *object.Boolean:
		r, ok := right.(*object.Boolean)
		if !ok {
			return newError(tok, "type mismatch: %s %s %s", left.Type(), op, right.Type())
		}
		eq = (l.Value == r.Value)
	case *object.String:
		r, ok := right.(*object.String)
		if !ok {
			return newError(tok, "type mismatch: %s %s %s", left.Type(), op, right.Type())
		}
		eq = (l.Value == r.Value)
	case *object.Null:
		_, ok := right.(*object.Null)
		eq = ok
	default:
		return newError(tok, "unknown equality operator comparison for %s", left.Type())
	}

	if op == "==" {
		return nativeBoolToBooleanObject(eq)
	}
	return nativeBoolToBooleanObject(!eq)
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

		res := Eval(we.Body, env)
		if res != nil {
			if res.Type() == object.RETURN_VALUE_OBJ || res.Type() == object.ERROR_OBJ {
				return res
			}
			result = res
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

		res := Eval(fe.Body, forEnv)
		if res != nil {
			if res.Type() == object.RETURN_VALUE_OBJ || res.Type() == object.ERROR_OBJ {
				return res
			}
			result = res
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

func applyFunction(fn object.Object, args []object.Object, tok token.Token) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		if len(args) != len(f.Parameters) {
			return newError(tok, "wrong number of arguments. got=%d, want=%d", len(args), len(f.Parameters))
		}
		extendedEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		res := f.Fn(args...)
		if err, ok := res.(*object.Error); ok {
			if err.Token.Line == 0 {
				err.Token = tok
			}
			return err
		}
		return res

	default:
		return newError(tok, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Name, args[i])
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

func newError(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Token:   tok,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
