package evaluator

import(
	"banana/ast"
	"banana/object"
)

var(
	NULL = &object.Null{}
	TRUE = &object.Boolean{Val: true}
	FALSE = &object.Boolean{Val: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Program:
		return evalStatements(node.Statements)
	// Expressions
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Val)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Op, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Val: node.Val}
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Op, right)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var res object.Object
	for _, stmt := range stmts {
		res = Eval(stmt)
	}
	return res
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case op == "==":
		return nativeBoolToBooleanObject(left == right)
	case op == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL	
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
		return NULL
	}
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
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
		return NULL
	}
	val := right.(*object.Integer).Val
	return &object.Integer{Val: -val}
}
