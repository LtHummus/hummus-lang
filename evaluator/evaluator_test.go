package evaluator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hummus-lang/lexer"
	"hummus-lang/object"
	"hummus-lang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-15", -15},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 *2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 * -10", -200},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestInnerReturnStatement(t *testing.T) {
	input := `
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5 * 5; let b = a; b;", 25},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	require.Truef(t, ok, "object is not a function. got %T (%+v)", evaluated, evaluated)

	require.Len(t, fn.Parameters, 1)

	assert.Equal(t, "x", fn.Parameters[0].String())
	assert.Equal(t, "(x + 2)", fn.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(add(5, 5), add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
let adder = fn(x) {
  fn(y) { x + y };
};

let addTwo = adder(2);
addTwo(2);
`
	testIntegerObject(t, testEval(input), 4)
}

func TestErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: can not + `INTEGER` and `BOOLEAN`",
		},
		{
			"5 + true; 5;",
			"type mismatch: can not + `INTEGER` and `BOOLEAN`",
		},
		{
			"-true",
			"unknown operator: unary - not defined for `BOOLEAN`",
		},
		{
			"true + false;",
			"unknown operator: binary + not defined for `BOOLEAN` and `BOOLEAN`",
		},
		{
			"5; true + false; 5",
			"unknown operator: binary + not defined for `BOOLEAN` and `BOOLEAN`",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: binary + not defined for `BOOLEAN` and `BOOLEAN`",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }
  return 1;
}
`,
			"unknown operator: binary + not defined for `BOOLEAN` and `BOOLEAN`",
		},
		{
			"foobar",
			"unknown reference: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !assert.Truef(t, ok, "error object not returned. got %T (%+v)", evaluated, evaluated) {
			continue
		}

		assert.Equal(t, tt.expected, errObj.Message)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) {
	assert.Equalf(t, Null, obj, "object is not null. Got %T (%+v)", obj, obj)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	if !assert.Truef(t, ok, "object is not an Boolean. Got %T (%+v)", obj, obj) {
		return
	}

	assert.Equal(t, expected, result.Value)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !assert.Truef(t, ok, "object is not an Integer. Got %T (%+v)", obj, obj) {
		return
	}

	assert.Equal(t, expected, result.Value)
}
