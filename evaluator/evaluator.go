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
	case *ast.Program:
		return evalProgram(node)
	// Statements
	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnVal)
		return &object.ReturnValue{Val: val}
	// Expressions
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Val)
	case *ast.IfExpression:
		return evalIfExpression(node)
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

func evalProgram(program *ast.Program) object.Object {
	var res object.Object
	for _, stmt := range program.Statements {
		res = Eval(stmt)
		if returnVal, ok := res.(*object.ReturnValue); ok {
			return returnVal.Val
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var res object.Object
	for _, stmt := range block.Statements {
		res = Eval(stmt)
		if res != nil && res.Type() == object.RETURN_VALUE_OBJ {
			return res
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
