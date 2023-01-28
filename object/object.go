package object

import (
	"bytes"
	"fmt"
	"strings"
)

type Object interface {
	Inspect() string
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

// ------------------------------
// Nill Object
// ------------------------------
type Nill struct {
}

func (n *Nill) Inspect() string {
	return "nil"
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

// ------------------------------
// Float Object
// ------------------------------
type Float struct {
	Value float64
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// ------------------------------
// Boolean Object
// ------------------------------
type Boolean struct {
	Value bool
}

func (f *Boolean) Inspect() string {
	return fmt.Sprintf("%t", f.Value)
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
