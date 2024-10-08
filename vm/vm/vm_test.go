package vm

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"testing"
	"vm/compiler"
)

type vmTestCase struct {
	input    string
	expected any
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50+100+-50", 0},
		{"(5+10*2+15/3)*2+-10", 50},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!!5", false},
		{"!5", true},
		{"!(if(false){5;})", true},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if(true){10}", 10},
		{"if(true){10}else{20}", 10},
		{"if(false){10}else{20}", 20},
		{"if(1){10}", 10},
		{"if(1<2){10}", 10},
		{"if(1<2){10}else{20}", 10},
		{"if(1>2){10}else{20}", 20},
		{"if(1>2){10}", Null},
		{"if(false){10}", Null},
	}
	runVmTests(t, tests)
}
func TestLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{`
    let one=1;
    let two=2;
    two;
    `, 2},
		{`
    let one=1;
    let two=2;
    one;
    `, 1},
		{
			`
      let one=1;
      let two=2;
      let three=one+two;
      three;`,
			3,
		},
	}
	runVmTests(t, tests)
}

func TestStringLiteral(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey";`, "monkey"},
		{`"mon"+"key";`, "monkey"},
	}
	runVmTests(t, tests)
}

func TestArrayLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []any{}},
		{"[1,2,3]", []any{1, 2, 3}},
		{"[1+2,2-3,3*4]", []any{3, -1, 12}},
		{`[1+2,2-3,3*4,"Hello","World"]`, []any{3, -1, 12, "Hello", "World"}},
	}
	runVmTests(t, tests)
}
func TestHashLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]any{}},
		{"{1:2,2:3}", map[object.HashKey]any{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 2}).HashKey(): 3,
		}},
		{"{1+2:2-3,2*3:3*4}", map[object.HashKey]any{
			(&object.Integer{Value: 3}).HashKey(): -1,
			(&object.Integer{Value: 6}).HashKey(): 12,
		}},
	}
	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1,2,3][1]", 2},
		{"[1,2,3][0+2]", 3},
		{"[[1,1,1]][0][0]", 1},
		{"[][0]", Null},
		{"[1,2,3][99]", Null},
		{"{1:1,2:2}[1]", 1},
		{"{1:1,2:2}[2]", 2},
		{"{1:1}[0]", Null},
		{"{}[0]", Null},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let fivePlusTen=fn(){5+10;};
			fivePlusTen()`,
			expected: 15,
		},
		{
			input: `
			let one=fn(){1;};
			let two=fn(){2;};
			one()+two();`,
			expected: 3,
		},
		{
			input: `
			let a = fn(){ 1 };
			let b = fn(){ a() + 1; };
			let c = fn(){ b() + 1; };
			c();
			`,
			expected: 3,
		},
	}
	runVmTests(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let earlyExit=fn(){return 99;100;}
			earlyExit()`,
			expected: 99,
		},
		{
			input: `
			let earlyExit=fn(){return 99;return 100;}
			earlyExit()`,
			expected: 99,
		},
	}
	runVmTests(t, tests)
}

func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let noReturn=fn(){};
			noReturn();`,
			expected: Null,
		},
		{
			input: `
			let noReturn=fn(){};
			let noReturnTwo=fn(){noReturn();};
			noReturn();
			noReturnTwo();`,
			expected: Null,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let one = fn() { let one = 1; one };
        one();
        `,
			expected: 1,
		},
		{
			input: `
        let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
        oneAndTwo();
        `,
			expected: 3,
		},
		{
			input: `
        let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
        let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
        oneAndTwo() + threeAndFour();
        `,
			expected: 10,
		},
		{
			input: `
        let firstFoobar = fn() { let foobar = 50; foobar; };
        let secondFoobar = fn() { let foobar = 100; foobar; };
        firstFoobar() + secondFoobar();
        `,
			expected: 150,
		},
		{
			input: `
        let globalSeed = 50;
        let minusOne = fn() {
            let num = 1;
            globalSeed - num;
        }
        let minusTwo = fn() {
            let num = 2;
            globalSeed - num;
        }
        minusOne() + minusTwo();
        `,
			expected: 97,
		},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
        let identity = fn(a) { a; };
        identity(4);
        `,
			expected: 4,
		},
		{
			input: `
        let sum = fn(a, b) { a + b; };
        sum(1, 2);
        `,
			expected: 3,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            c;
        };
        sum(1, 2);
        `,
			expected: 3,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            c;
        };
        sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input: `
        let sum = fn(a, b) {
            let c = a + b;
            c;
        };
        let outer = fn() {
            sum(1, 2) + sum(3, 4);
        };
        outer();
        `,
			expected: 10,
		},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for idx, tt := range tests {
		program := parse(tt.input)
		comp := compiler.NewCompiler()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error:%s", err)
		}
		vm := NewVM(comp.ByteCode())
		err = vm.Run()
		if err != nil {
			println("faild testcase ", idx)
			println(comp.ByteCode().Instructions.String())
			t.Fatalf("vm error: %s", err)
		}
		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}

}

func testExpectedObject(t *testing.T, expected any, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(bool(expected), actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []any:
		err := testArrayObject(expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[object.HashKey]any:
		err := testHashObject(t, expected, actual)
		if err != nil {
			t.Errorf("testHashObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null:%T %+v)", actual, actual)
		}

	}
}
func testHashObject(t *testing.T, expected map[object.HashKey]any, actual object.Object) error {
	result, ok := actual.(*object.HashObject)
	if !ok {
		return fmt.Errorf("object is not Hash. got=%T (%+v)", actual, actual)
	}
	if len(expected) != len(result.Pairs) {
		return fmt.Errorf("hash length is wrong. want=%d, got=%d", len(expected), len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			return fmt.Errorf("no pair for given key in Pairs")
		}
		testExpectedObject(t, expectedValue, pair.Value)
	}
	return nil
}

func testArrayObject(expecteds []any, actual object.Object) error {
	arr, ok := actual.(*object.ArrayObject)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%+v)", actual, actual)
	}
	if len(expecteds) != len(arr.Elements) {
		return fmt.Errorf("array length is wrong. want=%d, got=%d", len(expecteds), len(arr.Elements))
	}
	for i, o := range arr.Elements {
		switch expected := expecteds[i].(type) {
		case int64:
			err := testIntegerObject(expected, o)
			if err != nil {
				return err
			}
		case string:
			err := testStringObject(expected, o)
			if err != nil {
				fmt.Printf("%#v\n", arr.Elements)
				return err
			}
		}
	}
	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.StringObject)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. want=%q, got=%q", expected, result.Value)
	}
	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. want=%t, got=%t", expected, result.Value)
	}
	return nil
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

func parse(input string) *ast.Program {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return p.ParseProgram()
}
