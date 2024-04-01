package ast

import(
	"reflect"
	"testing"
)

func TestModify(t * testing.T) {
	one := func() Expression { return &IntegerLiteral{Val: 1} }
	two := func() Expression { return &IntegerLiteral{Val: 2 }}
	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}
		if integer.Val != 1 {
			return node
		}
		integer.Val = 2
		return integer
	}
	tests := []struct {
		input Node
		expected Node
	} {
		{
            one(),
            two(),
        },
        {
            &Program{
                Statements: []Statement{
                    &ExpressionStatement{Expression: one()},
                },
            },
            &Program{
                Statements: []Statement{
                    &ExpressionStatement{Expression: two()},
                },
            },
        },
		{
            &InfixExpression{Left: one(), Op: "+", Right: two()},
            &InfixExpression{Left: two(), Op: "+", Right: two()},
        },
        {
            &InfixExpression{Left: two(), Op: "+", Right: one()},
            &InfixExpression{Left: two(), Op: "+", Right: two()},
        },
		{
            &PrefixExpression{Op: "-", Right: one()},
            &PrefixExpression{Op: "-", Right: two()},
        },
		{
            &IndexExpression{Left: one(), Index: one()},
            &IndexExpression{Left: two(), Index: two()},
        },
		{
            &IfExpression{
                Condition: one(),
                Consequence: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: one()},
                    },
                },
                Alternative: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: one()},
                    },
                },
            },
            &IfExpression{
                Condition: two(),
                Consequence: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: two()},
                    },
                },
                Alternative: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: two()},
                    },
                },
            },
        },
		{
            &ReturnStatement{ReturnVal: one()},
            &ReturnStatement{ReturnVal: two()},
        },
		{
            &LetStatement{Val: one()},
            &LetStatement{Val: two()},
        },
		{
            &FunctionLiteral{
                Parameters: []*Identifier{},
                Body: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: one()},
                    },
                },
            },
            &FunctionLiteral{
                Parameters: []*Identifier{},
                Body: &BlockStatement{
                    Statements: []Statement{
                        &ExpressionStatement{Expression: two()},
                    },
                },
            },
        },
		{
            &ArrayLiteral{Elements: []Expression{one(), one()}},
            &ArrayLiteral{Elements: []Expression{two(), two()}},
        },
    }

    for _, tt := range tests {
        modified := Modify(tt.input, turnOneIntoTwo)

        equal := reflect.DeepEqual(modified, tt.expected)
        if !equal {
            t.Errorf("not equal. got=%#v, want=%#v",
                modified, tt.expected)
        }
	}

	dictLiteral := &DictLiteral{
        Pairs: map[Expression]Expression{
            one(): one(),
            one(): one(),
        },
    }

    Modify(dictLiteral, turnOneIntoTwo)

    for key, val := range dictLiteral.Pairs {
        key, _ := key.(*IntegerLiteral)
        if key.Val != 2 {
            t.Errorf("value is not %d, got=%d", 2, key.Val)
        }
        val, _ := val.(*IntegerLiteral)
        if val.Val != 2 {
            t.Errorf("value is not %d, got=%d", 2, val.Val)
        }
    }
}