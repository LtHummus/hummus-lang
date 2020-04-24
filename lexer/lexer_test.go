package lexer

import (
	"testing"

	"hummus-lang/token"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func TestSpecialCharacters(t *testing.T) {
	input := `=+(){},;*-/!<>`

	tests := []TestCase{
		{token.Assign, "="},
		{token.Plus, "+"},
		{token.LeftParen, "("},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.RightBrace, "}"},
		{token.Comma, ","},
		{token.Semicolon, ";"},
		{token.Asterisk, "*"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Bang, "!"},
		{token.Lt, "<"},
		{token.Gt, ">"},
		{token.Eof, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		assert.Equalf(t, tt.expectedType, tok.Type, "test %d failed (incorrect type)", i)
		assert.Equalf(t, tt.expectedLiteral, tok.Literal, "test %d failed (incorrect literal)", i)
	}
}

func TestRealCode(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);`

	tests := []TestCase{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},


		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LeftParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Semicolon, ";"},

		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LeftParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RightParen, ")"},
		{token.Semicolon, ";"},

		{token.Eof, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		assert.Equalf(t, tt.expectedType, tok.Type, "test %d failed (incorrect type)", i)
		assert.Equalf(t, tt.expectedLiteral, tok.Literal, "test %d failed (incorrect literal)", i)
	}

}

func TestAllKeywords(t *testing.T) {
	input := `if let return else true false`

	tests := []TestCase{
		{token.If, "if"},
		{token.Let, "let"},
		{token.Return, "return"},
		{token.Else, "else"},
		{token.True, "true"},
		{token.False, "false"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		assert.Equalf(t, tt.expectedType, tok.Type, "test %d failed (incorrect type)", i)
		assert.Equalf(t, tt.expectedLiteral, tok.Literal, "test %d failed (incorrect literal)", i)
	}
}

func TestMultiCharOperators(t *testing.T) {
	input := `if x != 5 == 4`
	tests := []TestCase{
		{token.If, "if"},
		{token.Ident, "x"},
		{token.NotEq, "!="},
		{token.Int, "5"},
		{token.Eq, "=="},
		{token.Int, "4"},
		{token.Eof, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		assert.Equalf(t, tt.expectedType, tok.Type, "test %d failed (incorrect type)", i)
		assert.Equalf(t, tt.expectedLiteral, tok.Literal, "test %d failed (incorrect literal)", i)
	}
}

func TestEverything(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
  x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"hello"
`

	tests := []TestCase{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LeftParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LeftParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RightParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
		{token.Gt, ">"},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LeftParen, "("},
		{token.Int, "5"},
		{token.Lt, "<"},
		{token.Int, "10"},
		{token.RightParen, ")"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Else, "else"},
		{token.LeftBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RightBrace, "}"},
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.String, "hello"},
		{token.Eof, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		assert.Equalf(t, tt.expectedType, tok.Type, "test %d failed (incorrect type)", i)
		assert.Equalf(t, tt.expectedLiteral, tok.Literal, "test %d failed (incorrect literal)", i)
	}

}
