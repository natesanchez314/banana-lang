package ast

import (
	"banana/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program {
		Statements: []Statement {
			&LetStatement {
				Token: token.Token {
					Type: token.LET, Literal: "let",
				},
				Name: &Identifier {
					Token: token.Token {
						Type: token.ID, Literal: "myVar",
					},
					Val: "myVar",
				},
				Val: &Identifier {
					Token: token.Token {
						Type: token.ID, Literal: "anotherVar",
					},
					Val: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}