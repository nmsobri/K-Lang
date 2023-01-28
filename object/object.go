package object

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	OBJECT_NILL    = "OBJECT_NILL"
	OBJECT_STRING  = "OBJECT_STRING"
	OBJECT_INTEGER = "OBJECT_INTEGER"
	OBJECT_FLOAT   = "OBJECT_FLOAT"
	OBJECT_BOOLEAN = "OBJECT_BOOLEAN"
	OBJECT_ARRAY   = "OBJECT_ARRAY"
	OBJECT_HASHMAP = "OBJECT_HASHMAP"
)

type Object interface {
	Inspect() string
	Type() ObjectType
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
	Value map[Object]Object
}

func (h *HashMap) Inspect() string {
	var out bytes.Buffer

	arrStr := []string{}

	for key, val := range h.Value {
		arrStr = append(arrStr, key.Inspect()+":"+val.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(arrStr, ", "))
	out.WriteString("}")
	return out.String()
}

func (h *HashMap) Type() ObjectType {
	return OBJECT_HASHMAP
}
