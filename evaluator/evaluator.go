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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	// Statements
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnVal)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Val: val}
	// Expressions
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Val)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Op, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Val: node.Val}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Op, right)
	}
	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Msg: fmt.Sprintf(format, a...)}
}

func evalProgram(program *ast.Program) object.Object {
	var res object.Object
	for _, stmt := range program.Statements {
		res = Eval(stmt)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Val
		case *object.Error:
			return res
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var res object.Object
	for _, stmt := range block.Statements {
		res = Eval(stmt)
		if res != nil {
			resType := res.Type()
			if resType == object.RETURN_VALUE_OBJ || resType == object.ERROR_OBJ {
				return res
			}
		}
	}
	return res
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt)
		if returnVal, ok := res.(*object.ReturnValue); ok {
			return returnVal.Val
		}
	}
	return res
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition) 
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("Type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), op, right.Type())
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
		return newError("Unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
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
