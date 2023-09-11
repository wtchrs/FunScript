package object

import (
	"bytes"
	"fmt"
	"funscript/ast"
	"hash/fnv"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	BUILTIN_OBJ      = "BUILTIN"
	ERROR_OBJ        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

var hashCache = make(map[Hashable]HashKey)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	if hashKey, ok := hashCache[i]; ok {
		return hashKey
	}

	hashKey := HashKey{Type: i.Type(), Value: uint64(i.Value)}
	hashCache[i] = hashKey
	return hashKey
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	if hashKey, ok := hashCache[s]; ok {
		return hashKey
	}

	h := fnv.New64a()
	h.Write([]byte(s.Value))
	hashKey := HashKey{Type: s.Type(), Value: h.Sum64()}

	hashCache[s] = hashKey
	return hashKey
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if hashKey, ok := hashCache[b]; ok {
		return hashKey
	}

	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	hashKey := HashKey{Type: b.Type(), Value: value}
	hashCache[b] = hashKey
	return hashKey
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var elements []string

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// TODO: Add stacktrace to Error

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
