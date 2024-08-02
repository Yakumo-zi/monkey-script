package compiler

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
	"vm/code"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "1+2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1-2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1/2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1*2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)

}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1>2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1<2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1==2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "1!=2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      if(true){10};3333;
      `,
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpTrue),
				//0001
				code.Make(code.OpJumpNotTruthy, 10),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpJump, 11),
				//0010
				code.Make(code.OpNull),
				//0011
				code.Make(code.OpPop),
				//0012
				code.Make(code.OpConstant, 1),
				//0015
				code.Make(code.OpPop),
			},
		},
		{
			input: `
      if(true){10;20};3333;
      `,
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpTrue),
				//0001
				code.Make(code.OpJumpNotTruthy, 14),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpPop),
				//0008
				code.Make(code.OpConstant, 1),
				//0011
				code.Make(code.OpJump, 15),
				//0014
				code.Make(code.OpNull),
				//0015
				code.Make(code.OpPop),
				//0016
				code.Make(code.OpConstant, 2),
				//0018
				code.Make(code.OpPop),
			},
		},
		{
			input: `
      if(true){10}else{20};3333;
      `,
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpTrue),
				//0001
				code.Make(code.OpJumpNotTruthy, 10),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpJump, 13),
				//00010
				code.Make(code.OpConstant, 1),
				//0013
				code.Make(code.OpPop),
				//0014
				code.Make(code.OpConstant, 2),
				//0017
				code.Make(code.OpPop),
			},
		},
		{
			input: `
      if(false){10};3333;
      `,
			expectedConstants: []any{10, 3333},
			expectedInstructions: []code.Instructions{
				//0000
				code.Make(code.OpFalse),
				//0001
				code.Make(code.OpJumpNotTruthy, 10),
				//0004
				code.Make(code.OpConstant, 0),
				//0007
				code.Make(code.OpJump, 11),
				//00010
				code.Make(code.OpNull),
				//0011
				code.Make(code.OpPop),
				//0012
				code.Make(code.OpConstant, 1),
				//0015
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
      let one=1;
      let tow=2`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
      let one=1;
      one;`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		}, {
			input: `
      let one=1;
      let two=one;
      two;`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"monkey"`,
			expectedConstants: []any{"monkey"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input:             `"mon"+"key"`,
			expectedConstants: []any{"mon", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for idx, tt := range tests {
		program := parse(tt.input)
		compiler := NewCompiler()
		err := compiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error:%s", err)
		}
		bytecode := compiler.ByteCode()
		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions error:%s\nidx=%d", err, idx)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants error:%s\nidx=%d", err, idx)
		}
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. want=%d, got=%d", expected, result.Value)
	}
	return nil
}

func testConstants(t *testing.T, expected []any, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. want=%d, got=%d", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed:%s", i, err)
			}
		}
	}
	return nil
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length,\nwant=%q\ngot=%q", concatted, actual)
	}
	for i, ins := range concatted {
		if actual[i] != ins {
			fmt.Println("expected:")
			fmt.Println(concatted.String())
			fmt.Println("actual:")
			fmt.Println(actual.String())
			return fmt.Errorf("wrong instruction at %d.\nwant=%q\ngot=%q", i, ins, actual[i])
		}
	}
	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}
	return out
}

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.ParseProgram()
}
