package object

import (
	"bytes"
	"fmt"
	"hummus-lang/ast"
	"strings"
)

type ObjectType string
type PredefFunction func(args ...Object) Object

const (
	IntegerObj         = "INTEGER"
	BooleanObj         = "BOOLEAN"
	StringObj          = "STRING"
	NullObj            = "NULL"
	ReturnValueObj     = "RETURN_VALUE"
	ErrorObj           = "ERROR"
	FunctionObj        = "FUNCTION"
	PredefinedFunction = "PREDEFINED_FUNCTION"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Printable() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string   { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType  { return IntegerObj }
func (i *Integer) Printable() string { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Inspect() string   { return fmt.Sprintf(`"%s"`, s.Value) }
func (s *String) Type() ObjectType  { return StringObj }
func (s *String) Printable() string { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Printable() string {
	if b.Value {
		return "true"
	} else {
		return "false"
	}
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string   { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType  { return ReturnValueObj }
func (rv *ReturnValue) Printable() string { return "RV" }

type Null struct{}

func (n *Null) Inspect() string   { return "null" }
func (n *Null) Type() ObjectType  { return NullObj }
func (n *Null) Printable() string { return "null" }

type Error struct {
	Message string
}

func (e *Error) Inspect() string   { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType  { return ErrorObj }
func (e *Error) Printable() string { return fmt.Sprintf("error: %s", e.Message) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Type() ObjectType  { return FunctionObj }
func (f *Function) Printable() string { return "user defined function" } //TODO: change this

type Predef struct {
	Function PredefFunction
}

func (f *Predef) Inspect() string   { return "predefined function" }
func (f *Predef) Type() ObjectType  { return PredefinedFunction }
func (f *Predef) Printable() string { return "predefined function" }
