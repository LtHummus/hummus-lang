package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hummus-lang/ast"
	"hummus-lang/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo = 83838;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.NotNil(t, program, "parser returned nil")

	require.Lenf(t, program.Statements, 3, "3 statements should have been returned")

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
	}
}

func TestInvalidStatement(t *testing.T) {
	input := `
let x 5;
let = 10;
let 123455
`
	l := lexer.New(input)
	p := New(l)
	_ = p.ParseProgram()

	errors := p.Errors()
	require.Len(t, errors, 4)

}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 98765;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	require.Empty(t, p.Errors())
	require.Len(t, program.Statements, 3)

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		assert.True(t, ok)

		assert.Equal(t, "return", returnStmt.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.Errors())
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Statement is not ExpressionStatement. Got %T instead", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Identifier)
	require.Truef(t, ok, "Expression not *ast.Identifier. Got %T instead", ident.Value)

	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLiteral())
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.Errors())
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Statement is not ExpressionStatement. Got %T instead", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	require.Truef(t, ok, "Expression not *ast.IntegerLiteral. Got %T instead", ident)

	assert.Equal(t, int64(5), ident.Value)
	assert.Equal(t, "5", ident.TokenLiteral())
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Statement is not ExpressionStatement. Got %T instead", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Boolean)
	require.Truef(t, ok, "Expression not *ast.Boolean. Got %T instead", ident)

	assert.Equal(t, true, ident.Value)
}

func TestBooleanLiteralExpressionFalse(t *testing.T) {
	input := "false;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Statement is not ExpressionStatement. Got %T instead", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Boolean)
	require.Truef(t, ok, "Expression not *ast.Boolean. Got %T instead", ident)

	assert.Equal(t, false, ident.Value)
}

func TestParsingPrefixExpressionBang(t *testing.T) {
	input := "!5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.Errors(), "program did not parse without errors")
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "statement not ast.ExpressionStatement. Got %T instead", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	require.Truef(t, ok, "expression not ast.PrefixExpression. Got %T instead", stmt.Expression)
	require.Equal(t, "!", exp.Operator)

	testIntegerLiteral(t, exp.Right, 5)
}

func TestParsingPrefixExpressionMinus(t *testing.T) {
	input := "-5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.Errors(), "program did not parse without errors")
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "statement not ast.ExpressionStatement. Got %T instead", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.PrefixExpression)
	require.Truef(t, ok, "expression not ast.PrefixExpression. Got %T instead", stmt.Expression)
	require.Equal(t, "-", exp.Operator)

	testIntegerLiteral(t, exp.Right, 5)
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		require.Empty(t, p.errors)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.Truef(t, ok, "statement is not an ast.ExpressionStatement. got %T instead", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		require.Truef(t, ok, "expression is not ast.InfixExpression. got %T instead", stmt.Expression)

		testIntegerLiteral(t, exp.Left, tt.leftValue)

		require.Equal(t, tt.operator, exp.Operator)

		testIntegerLiteral(t, exp.Right, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b * c", "(a + (b * c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
	}

	for idx, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		require.Emptyf(t, p.errors, "test case %d failed", idx)

		actual := program.String()
		assert.Equalf(t, tt.expected, actual, "test case %d failed", idx)
	}
}

func TestOperatorPrecedenceParsingWithParens(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
	}

	for idx, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		require.Emptyf(t, p.errors, "test case %d failed", idx)

		actual := program.String()
		assert.Equalf(t, tt.expected, actual, "test case %d failed", idx)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "statement is not ast.ExpressionStatement. Got %T instead", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expression is not IfExpression. Got %T instead", stmt.Expression)

	assert.Len(t, exp.Consequence.Statements, 1)

	conditional, ok := exp.Condition.(*ast.InfixExpression)
	require.Truef(t, ok, "condition not infix expression. got %T instead", exp.Condition)

	require.Equal(t, "<", conditional.Operator)
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "statement is not ast.ExpressionStatement. Got %T instead", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.Truef(t, ok, "expression is not IfExpression. Got %T instead", stmt.Expression)

	assert.Len(t, exp.Consequence.Statements, 1)

	conditional, ok := exp.Condition.(*ast.InfixExpression)
	require.Truef(t, ok, "condition not infix expression. got %T instead", exp.Condition)

	require.Equal(t, "<", conditional.Operator)

}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "statement is not an ast.ExpressionStatement. Got %T instead", program.Statements[0])

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.Truef(t, ok, "expression was not a FunctionLiteral. Got %T instead", stmt.Expression)

	assert.Len(t, function.Parameters, 2)
	assert.Equal(t, function.Parameters[0].Value, "x")
	assert.Equal(t, function.Parameters[1].Value, "y")

	assert.Len(t, function.Body.Statements, 1)

	_, ok = function.Body.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "body was not an expressionStatement. Got %T instead", function.Body.Statements[0])


}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	require.Empty(t, p.errors)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.Truef(t, ok, "Expected ExpressionStatement, got %T instead", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	require.Truef(t, ok, "Expected CallExpression, got %T instead", stmt.Expression)

	assert.Len(t, exp.Arguments, 3)
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	if !assert.Truef(t, ok, "il not *ast.IntegerLiteral. got %T instead", il) {
		return
	}

	assert.Equal(t, value, integ.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), integ.TokenLiteral())
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "let", s.TokenLiteral())

	letStmt, ok := s.(*ast.LetStatement)
	require.Truef(t, ok, "not *ast.LetStatement. got %T", s)

	require.Equalf(t, name, letStmt.Name.TokenLiteral(), "incorrect name")
}
