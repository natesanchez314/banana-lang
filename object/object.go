package object

import(
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"banana/ast"
)

type ObjectType string

const (
	ARRAY_OBJ = "ARRAY"
	BOOLEAN_OBJ = "BOOLEAN"
	BUILTIN_OBJ = "BUILTIN"
	DICT_OBJ = "DICT"
	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	INTEGER_OBJ = "INTEGER"
	MACRO_OBJ = "MACRO"
	NULL_OBJ = "NULL"
	QUOTE_OBJ = "QUOTE"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	STRING_OBJ = "STRING"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	DictKey() DictKey
}

type Array struct {
	Elements []Object
}
func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Boolean struct {
	Val bool
}
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Val) }
func (b *Boolean) DictKey() DictKey {
	var val uint64
	if b.Val {
		val = 1
	} else {
		val = 0
	}
	return DictKey{Type: b.Type(), Val: val}
}

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Fn BuiltinFunction
}
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string { return "builtin function" }

type Dict struct {
	Pairs map[DictKey]DictPair
}
func (d *Dict) Type() ObjectType { return DICT_OBJ }
func (d *Dict) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range d.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Val.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
} 

type DictKey struct {
	Type ObjectType
	Val uint64
}

type DictPair struct {
	Key Object
	Val Object
}

type Error struct {
	Msg string
}
func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Msg }

type Function struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Environment
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
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

type Integer struct {
	Val int64
}
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Val) }
func (i *Integer) DictKey() DictKey { return DictKey{Type: i.Type(), Val: uint64(i.Val)} }

type Macro struct {
	Parameters []*ast.Identifier
	Body *ast.BlockStatement
	Env *Environment
}
func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Null struct {}
func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

type Quote struct {
	Node ast.Node
}
func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string { return "QUOTE(" + q.Node.String() + ")" }

type ReturnValue struct {
	Val Object
}
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Val.Inspect() }

type String struct {
	Val string
}
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string { return s.Val }
func (s *String) DictKey() DictKey {
	h := fnv.New64a()
	h.Write([]byte(s.Val))
	return DictKey{Type: s.Type(), Val: h.Sum64()}
}