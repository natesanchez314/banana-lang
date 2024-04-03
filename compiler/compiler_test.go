package compiler

import (
	"testing"
	"banana/code"
)

type compilerTestCase struct {
	input string
	expectedConstants []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase {
		{
			input: "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}
	runCompilerTests(t. tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()
	for _, tt := range tests {
		program := parse(tt.input)
		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compiler error: %s", err)
		}
		byteCode := compiler.ByteCode()
		err = testInstructions(tt.expectedInstructions, byteCode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}
		err = testConstants(t, tt.expectedConstants, byteCode.Constants)
		if err != nil {
			t.Fatalf("testConstans failed: %s", err)
		}
	}
}