package object

import (
	"Klang/ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	OBJECT_NILL     = "OBJECT_NILL"
	OBJECT_STRING   = "OBJECT_STRING"
	OBJECT_INTEGER  = "OBJECT_INTEGER"
	OBJECT_FLOAT    = "OBJECT_FLOAT"
	OBJECT_BOOLEAN  = "OBJECT_BOOLEAN"
	OBJECT_ARRAY    = "OBJECT_ARRAY"
	OBJECT_HASHMAP  = "OBJECT_HASHMAP"
	OBJECT_FUNCTION = "OBJECT_FUNCTION"
)

type Object interface {
	Inspect() string
	Type() ObjectType
}

type Hash struct {
	Type  ObjectType
	Value string
}

type Hashable interface {
	Hashkey() Hash
}

// ------------------------------
// String Object
// ------------------------------
type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return OBJECT_STRING
}

func (s *String) Hashkey() Hash {
	return Hash{Type: OBJECT_STRING, Value: s.Value}
}

// ------------------------------
// Nill Object
// ------------------------------
type Nill struct {
}

func (n *Nill) Inspect() string {
	return "nil"
}

func (n *Nill) Type() ObjectType {
	return OBJECT_NILL
}

// ------------------------------
// Integer Object
// ------------------------------
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return OBJECT_INTEGER
}

func (i *Integer) Hashkey() Hash {
	return Hash{Type: OBJECT_INTEGER, Value: fmt.Sprintf("%d", i.Value)}
}

// ------------------------------
// Float Object
// ------------------------------
type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

func (f *Float) Type() ObjectType {
	return OBJECT_FLOAT
}

// ------------------------------
// Boolean Object
// ------------------------------
type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return OBJECT_BOOLEAN
}

// ------------------------------
// Array Object
// ------------------------------
type Array struct {
	Value []Object
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	arrStr := []string{}

	for _, elem := range a.Value {
		arrStr = append(arrStr, elem.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(arrStr, ", "))
	out.WriteString("]")
	return out.String()
}

func (a *Array) Type() ObjectType {
	return OBJECT_ARRAY
}

// ------------------------------
// HashMap Object
// ------------------------------
type HashMap struct {
	Value map[Hash]Object
}

func (h *HashMap) Inspect() string {
	var out bytes.Buffer

	arrStr := []string{}

	for key, val := range h.Value {
		arrStr = append(arrStr, key.Value+":"+val.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(arrStr, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *HashMap) Type() ObjectType {
	return OBJECT_HASHMAP
}

// ------------------------------
// Function Object
// ------------------------------
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}

	for _, val := range f.Parameters {
		params = append(params, val.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

func (f *Function) Type() ObjectType {
	return OBJECT_FUNCTION
}
