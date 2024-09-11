package object

import (
	"bytes"
)

type Environment struct {
    store map[string]Object
    outer *Environment
}

func NewEnvironment() *Environment {
    s := make(map[string]Object)
    return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}

func (e *Environment) Get(key string) (Object, bool) {
    val, ok := e.store[key]
    if !ok && e.outer != nil {
        val, ok = e.outer.Get(key)
    }
    return val, ok
}

func (e *Environment) Set(key string, val Object) Object {
    e.store[key] = val
    return val
}

func (e *Environment) String() string {
    var out bytes.Buffer

    out.WriteString("store\n")
    for key, val := range e.store {
        out.WriteString("Key: ")
        out.WriteString(key)
        out.WriteString(",\t")
        out.WriteString("Value: ")
        out.WriteString(val.Inspect())
        out.WriteString("\n")
    }

    if e.outer != nil {
        out.WriteString("outer\n")
        for key, val := range e.outer.store {
            out.WriteString("Key: ")
            out.WriteString(key)
            out.WriteString(",\t")
            out.WriteString("Value: ")
            out.WriteString(val.Inspect())
            out.WriteString("\n")
        }
    }

    return out.String()
}
