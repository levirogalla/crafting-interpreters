package environ

import (
	lerr "crafting-interpreters/error"
	m "crafting-interpreters/models"
	"fmt"
)


type Environ struct {
	values map[string]m.Ltype
	enclosing *Environ
}

func NewEnviron(enclosing *Environ) *Environ {
	return &Environ{
		values: make(map[string]m.Ltype),
		enclosing: enclosing,
	}
}

func (e *Environ) Define(name string, value m.Ltype) {
	e.values[name] = value
}
	
func (e *Environ) Get(name m.Token) (m.Ltype, error) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Get(name)
		}
		return nil, lerr.RuntimeError{
			Token: name,
			Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
		}
	}
	return val, nil
}

func (e *Environ) Assign(name m.Token, value m.Ltype) error {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return lerr.RuntimeError{
		Token: name,
		Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
	}
}
