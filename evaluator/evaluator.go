package evaluator

import(
	"fmt"
	"banana/ast"
	"banana/object"
)

var(
	NULL = &object.Null{}
	TRUE = &object.Boolean{Val: true}
	FALSE = &object.Boolean{Val: false}
)

const(
	unknownInfixOp string = "Unknown operator: %s %s %s"
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	// Statements
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.LetStatement:
		val := Eval(node.Val, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Val, val)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnVal, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Val: val}
	// Expressions
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Val)
	case *ast.CallExpression:
		function := Eval(node.Fun, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Op, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Val: node.Val}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Op, right)
	case *ast.StringLiteral:
		return &object.String{Val: node.Val}
	}
	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range program.Statements {
		res = Eval(stmt, env)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Val
		case *object.Error:
			return res
		}
	}
	return res
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt, env)
		if returnVal, ok := res.(*object.ReturnValue); ok {
			return returnVal.Val
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range block.Statements {
		res = Eval(stmt, env)
		if res != nil {
			resType := res.Type()
			if resType == object.RETURN_VALUE_OBJ || resType == object.ERROR_OBJ {
				return res
			}
		}
	}
	return res
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var res []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("Not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Val, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Val
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(i.Val); ok {
		return val
	}
	if builtin, ok := builtins[i.Val]; ok {
		return builtin
	}
	return newError("Identifier not found: " + i.Val)
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

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left, right)
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError(unknownInfixOp, left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Val
	rightVal := right.(*object.Integer).Val
	switch op {
	case "+":
		return &object.Integer{Val: leftVal + rightVal}
	case "-":
		return &object.Integer{Val: leftVal - rightVal}
	case "*":
		return &object.Integer{Val: leftVal * rightVal}
	case "/":
		return &object.Integer{Val: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(unknownInfixOp, left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right object.Object) object.Object {
	if op != "+" {
		return newError("Unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
	leftVal := left.(*object.String).Val
	rightVal := right.(*object.String).Val
	return &object.String{Val: leftVal + rightVal}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator: %s%s", op, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}
	val := right.(*object.Integer).Val
	return &object.Integer{Val: -val}
}
