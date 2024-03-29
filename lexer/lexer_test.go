package lexer

import (
	"testing"
	"banana/token"
)

func TestNextToken(t *testing.T) {
	input := `
		let five = 5;
		"foobar"
		"foo bar"
		"[1, 2];"
		`

	tests := []struct {
		expectedType token.TokenType
		expectedLiteral string
	} {
		{token.LET, "let"},
		{token.ID, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. Expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}