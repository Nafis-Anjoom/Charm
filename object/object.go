package object

import (
	"bytes"
	"charm/ast"
	"fmt"
	"strings"
	"hash/fnv"
)

type ObjectType string


const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASHMAP_OBJ      = "HASHMAP"
	PAIR_OBJ         = "PAIR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type String struct {
	Value string
}
func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}

type Boolean struct {
	Value bool
}
func (i *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}

type Null struct {}
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return fmt.Sprintf("null")
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR:" + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
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
	out.WriteString("\n}\n")

	return out.String()
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, obj := range a.Elements {
		elements = append(elements, obj.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Pair struct {
	Key   Object
	Value Object
}
func (p *Pair) Type() ObjectType {
	return PAIR_OBJ
}
func (p *Pair) Inspect() string {
	return p.Key.Inspect() + ": " + p.Value.Inspect()
}

type HashMap struct {
	Map map[uint64]Pair
}
func (hm *HashMap) Type() ObjectType {
	return HASHMAP_OBJ
}
func (hm *HashMap) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	for _, pair := range hm.Map {
		pairs = append(pairs, pair.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// for debugging and testing purposes, Object interface is being nested.
// necessary for viewing the value through Inspect()
// TODO: investigate the performance impact of nested interfaces
type Hashable interface {
    Object
	HashCode() uint64
}
func (i *Integer) HashCode() uint64 {
	return uint64(i.Value)
}
func (b *Boolean) HashCode() uint64 {
	if b.Value {
		return 1
	}
	return 0
}
func (s *String) HashCode() uint64 {
	h := fnv.New64()
	h.Write([]byte(s.Value))
	return h.Sum64()
}
