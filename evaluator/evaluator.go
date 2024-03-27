package evaluator

import(
	"banana/ast"
	"banana/object"
)

var(
	NULL = &object.Boolean{}
	TRUE = &object.Boolean{Val: true}
	FALSE = &object.Boolean{Val: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Val)
	case *ast.IntegerLiteral:
		return &object.Integer{Val: node.Val}
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