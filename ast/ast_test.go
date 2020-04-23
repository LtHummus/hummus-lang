package ast

import (
	"github.com/stretchr/testify/assert"
	"hummus-lang/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement {
			&LetStatement{
				Token: token.Token{
					Type:    token.Let,
					Literal: "let",
					Line:    1,
				},
				Name:  &Identifier{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "myVar",
						Line:    1,
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.Ident,
						Literal: "anotherVar",
						Line:    1,
					},
					Value: "anotherVar",
				},
			},
		},
	}

	assert.Equal(t, "let myVar = anotherVar;", program.String())
}