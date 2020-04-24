package lexer

import "hummus-lang/token"

type Lexer struct {
	input        string
	position     int  // current position in input (current char)
	readPosition int  // current reading position (after current char)
	ch           byte // current char

	line int
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte, line int) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
		}
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		} else if l.ch == 0 {
			break //TODO: report error for non-ended string
		}
	}

	return l.input[pos:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Eq, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.Assign, l.ch, l.line)
		}
	case '+':
		tok = newToken(token.Plus, l.ch, l.line)
	case '-':
		tok = newToken(token.Minus, l.ch, l.line)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NotEq, Literal: literal, Line: l.line}
		} else {
			tok = newToken(token.Bang, l.ch, l.line)
		}
	case '"':
		literal := l.readString()
		tok = token.Token{Type: token.String, Literal: literal, Line: l.line}
	case '/':
		tok = newToken(token.Slash, l.ch, l.line)
	case '*':
		tok = newToken(token.Asterisk, l.ch, l.line)
	case '<':
		tok = newToken(token.Lt, l.ch, l.line)
	case '>':
		tok = newToken(token.Gt, l.ch, l.line)
	case ';':
		tok = newToken(token.Semicolon, l.ch, l.line)
	case '(':
		tok = newToken(token.LeftParen, l.ch, l.line)
	case ')':
		tok = newToken(token.RightParen, l.ch, l.line)
	case ',':
		tok = newToken(token.Comma, l.ch, l.line)
	case '{':
		tok = newToken(token.LeftBrace, l.ch, l.line)
	case '}':
		tok = newToken(token.RightBrace, l.ch, l.line)
	case 0:
		tok.Literal = ""
		tok.Type = token.Eof
		tok.Line = l.line
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = l.line
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.Int
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch, l.line)
		}
	}

	l.readChar()
	return tok
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}
