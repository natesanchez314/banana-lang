package parser

import (
	"testing"
	"banana/ast"
	"banana/lexer"
	"fmt"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Exression is not ast.CallExpressionm got-%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Fun, "add") {
		return
	}
	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of args, got=%d", len(exp.Args))
	}
	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input string
		expectedParams []string
	} {
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input) 
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fun := stmt.Expression.(*ast.FunctionLiteral)
		if len(fun.Parameters) != len(tt.expectedParams) {
			t.Errorf("length of parameters wrong. Expected %d, got=%d.\n", len(tt.expectedParams), len(fun.Parameters))
		}

		for i, id := range tt.expectedParams {
			testLiteralExpression(t, fun.Parameters[i], id)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x+ y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	fun, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Exression is not ast.FunctionLiteral, got =%T", stmt.Expression)
	}

	if len(fun.Parameters) != 2 {
		t.Fatalf("fun literal parameters wrong. Expected %d, got=%d\n", 2, len(fun.Parameters))
	}
	testLiteralExpression(t, fun.Parameters[0], "x")
	testLiteralExpression(t, fun.Parameters[1], "y")

	if len(fun.Body.Statements) != 1 {
		t.Fatalf("fun body stmt is not ast.ExpressionStatement, got=%T", fun.Body.Statements[0])
	}


	bodyStmt, ok := fun.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fun body stmt is not ast.ExpressionStatement, got=%T", fun.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Exression is not ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements, got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil, got=%+v", exp.Alternative)
	}
}

// func TestOperatorPrecedence(t *testing.T) {
// 	tests := []struct {
// 		input string
// 		expected string
// 	} {
// 		{
// 			"true",
// 			"true",
// 		},
// 		{
// 			"false",
// 			"false",
// 		},
// 		{
// 			"3 > 5 == false",
// 			"((3 > 5) == false)",
// 		},
// 		{
// 			"3 < 5 == true",
// 			"((3 < 5) == true)",
// 		},
// 	}

// 	for _, tt := range tests {
// 		l := lexer.New(tt.input)
// 		p := New(l)
// 		program := p.ParseProgram()
// 		checkParserErrors(t, p)

// 		if len(program.Statements) != 4  {
// 			t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 4, len(program.Statements))
// 		}
		
// 		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 		if !ok {
// 			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
// 		}

// 		exp, ok := stmt.Expression.(*ast.InfixExpression)
// 		if !ok {
// 			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
// 		}

// 		if exp.Op != tt.op {
// 			t.Fatalf("exp.Op is not '%s', got=%s", tt.op, exp.Op)
// 		}

// 		if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.op, tt.rightVal) {
// 			return
// 		}
// 	}
// }

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input string
		leftVal interface{}
		op string
		rightVal interface{}
	} {
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1  {
			t.Fatalf("program.Statements does not contain %d statements, got=%d\n", 1, len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Op != tt.op {
			t.Fatalf("exp.Op is not '%s', got=%s", tt.op, exp.Op)
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftVal, tt.op, tt.rightVal) {
			return
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		op string
		val interface{}
	} {
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if exp.Op != tt.op {
			t.Fatalf("exp.Op not ast.PrefixExpression, got=%T", stmt.Expression)
		}
		if !testLiteralExpression(t, exp.Right, tt.val) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%q", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got=%T", stmt.Expression)
	}
	if literal.Val != 5 {
		t.Errorf("literal.Val %d, got=%d", 5, literal.Val)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s, got=%s", "5", literal.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements, got=%q", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	id, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if id.Val != "foobar" {
		t.Errorf("id.Val not %s, got=%s", "foobar", id.Val)
	}
	if id.TokenLiteral() != "foobar" {
		t.Errorf("id.TokenLiteral not %s, got=%s", "foobar", id.TokenLiteral())
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil.")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 Statements, got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatment. Got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got%q", returnStmt.TokenLiteral())
		}
	}
}

func TestLetStatments(t *testing.T) {
	tests := []struct {
		input string
		expectedId string
		expectedVal interface{}
	} {
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests{
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements, got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testLetStatements(t, stmt, tt.expectedId) {
			return
		}
		val := stmt.(*ast.LetStatement).Val
		if !testLiteralExpression(t, val, tt.expectedVal) {
			return
		}
	}
}

func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. Got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. Got=%T", s)
		return false
	}

	if letStmt.Name.Val != name {
		t.Errorf("letStmt.Name.Value not '%s'. Got=%s", name, letStmt.Name.Val)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. Got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors.", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q.", msg)
	}
	t.FailNow()
}


func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}
	if integ.Val != val {
		t.Errorf("integ.Value not %d, got=%d", val, integ.Val)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integ.TokenLiteral not %d, got=%s", val, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, val string) bool {
	id, ok:= exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got=%T", exp)
		return false
	}
	if id.Val != val {
		t.Errorf("id.Val not %s, got=%s", val, id.Val)
		return false
	}
	if id.TokenLiteral() != val {
		t.Errorf("id.TokenLiteral not %s, got=%s", val, id.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, val bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean, got=%T", exp)
		return false
	}
	if b.Val != val {
		t.Errorf("b.Val not %t, got=%t", val, b.Val)
		return false
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", val) {
		t.Errorf("b.TokenLiteral not %t, got=%s", val, b.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Type of exp not handled, got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, op string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Op != op {
		t.Errorf("exp.Op is not '%s', got=%q", op, opExp.Op)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}