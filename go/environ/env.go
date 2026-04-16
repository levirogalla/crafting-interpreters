package environ

import (
	lerr "crafting-interpreters/error"
	m "crafting-interpreters/models"
	"fmt"
)


type Environ struct {
	values map[string]m.Ltype
}

func NewEnviron() *Environ {
	return &Environ{
		values: make(map[string]m.Ltype),
	}
}

func (e *Environ) Define(name string, value m.Ltype) {
	e.values[name] = value
}
	
func (e *Environ) Get(name m.Token) (m.Ltype, error) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		return nil, lerr.RuntimeError{
			Token: name,
			Message: fmt.Sprintf("undefined variable '%s'.", name.Lexeme),
		}
	}
	return val, nil
}