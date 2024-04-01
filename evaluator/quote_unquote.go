package evaluator

import(
	"fmt"
	"banana/ast"
	"banana/object"
	"banana/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		if len(call.Args) != 1 {
			return node
		}
		unquoted := Eval(call.Args[0], env)
		return convertObjectToAstNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Fun.TokenLiteral() == "unquote"
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Boolean:
		var t token.Token
		if obj.Val {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Val: obj.Val}
	case *object.Integer:
		t := token.Token{
			Type: token.INT,
			Literal: fmt.Sprintf("%d", obj.Val),
		}
		return &ast.IntegerLiteral{Token: t, Val: obj.Val}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}