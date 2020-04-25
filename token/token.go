package token

const (
	Illegal = "ILLEGAL"
	Eof     = "EOF"

	// identifiers & literals
	Ident  = "IDENT"
	Int    = "INT"
	String = "STRING"

	// Operators
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang     = "!"
	Asterisk = "*"
	Slash    = "/"
	Percent  = "%"

	Eq    = "=="
	NotEq = "!="

	Lt = "<"
	Gt = ">"

	// Delimiters
	Comma     = ","
	Semicolon = ";"

	LeftParen    = "("
	RightParen   = ")"
	LeftBrace    = "{"
	RightBrace   = "}"
	LeftBracket  = "["
	RightBracket = "]"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
