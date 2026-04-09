package interpreter

import (
	"crafting-interpreters/ast"
	lerr "crafting-interpreters/error"
	m "crafting-interpreters/models"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Interp struct{
	errReporter lerr.Reporter
}

func (i *Interp) Interpret(expr ast.Expr) {
	value, err := i.evaluate(expr)
	if err != nil {
		t :=(err.(*RuntimeError).token)
		m :=(err.(*RuntimeError).message)
		i.errReporter.RuntimeError(&t, &m)
		return
	} 
	fmt.Println(stringify(value))
}

func (*Interp) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value.Lit, nil
}

func (i *Interp) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interp) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	var zNum m.Lnum
	var zStr string
	left, err := i.evaluate(expr.Left)
	if err != nil { return nil, err }
	right, err := i.evaluate(expr.Right)
	if err != nil { return nil, err }

	switch expr.Op.TType {
	case m.Minus:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l - r, nil
	case m.Plus:
		switch l := left.(type) {
		case m.Lnum:
			switch r := right.(type) {
			case m.Lnum:
				return l + r, nil
			default:
				return nil, invalidOperandError(expr.Op, r, zNum)
			}
		case m.Lstring:
			switch r := right.(type) {
			case m.Lstring:
				return l + r, nil
			default:
				return nil, invalidOperandError(expr.Op, r, zStr)
			}
		default:
			return nil, invalidOperandError(expr.Op, l, zStr, zNum)
		}
	case m.Slash:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l / r, nil
	case m.Star:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l * r, nil
	case m.GT:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l > r, nil
	case m.GTE:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l >= r, nil
	case m.LT:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l < r, nil
	case m.LTE:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil { return nil, err }
		return l <= r, nil
	case m.Eq:
		return isEq(left, right), nil
	case m.Neq:
		return !isEq(left, right), nil
	default:
		panic("unreachable")
	}
}

func (i *Interp) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil { return nil, err }
	switch expr.Op.TType {
	case m.Bang:
		return !i.isTruthy(right), nil
	case m.Minus:
		r, err := checkType[m.Lnum](expr.Op, right)
		if err != nil { return nil, err }
		return -r, nil
	default:
		panic("unreachable: switch unary operator type")
	}
}

func (i *Interp) evaluate(expr ast.Expr) (any, error) {
	return ast.Accept(expr, i)
}

func (i *Interp) isTruthy(expr any) bool {
	if expr == nil {
		return false
	}
	switch e := expr.(type) {
	case bool:
		return e
	default:
		return true
	}
}

func isEq(l any, r any) bool {
	if l == nil && r == nil { 
		return true 
	} else if l == nil || r == nil {
		return false
	}

	switch l_ := l.(type) {
	case m.Lstring, m.Lnum: 
		switch r_ := r.(type) {
		case m.Lstring, m.Lnum:
			return l_ == r_
		}
		return false
	case []any:
		switch r_ := r.(type) {
		case []any:
			return slices.Equal(l_, r_)
		}
		return false
	default:
		return reflect.DeepEqual(l, r)
	}
}

func checkTypes[T any](op *m.Token, a, b any) (a_ T, b_ T, err error) {
	var ok bool
	var zero T
	a_, ok = a.(T)
	if !ok { 
		return zero, zero, invalidOperandError(op, a, zero)
	}
	b_, ok = b.(T)
	if !ok { 
		return zero, zero, invalidOperandError(op, a, zero)
 }
	return a_, b_, nil
}

func checkType[T any](op *m.Token, a any) (a_ T, err error) {
	var ok bool
	var zero T
	a_, ok = a.(T)
	if !ok { 
		return zero, invalidOperandError(op, a, zero)
	}
	return a_, nil
}

type RuntimeError struct {
	token m.Token
	message string
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", r.message, r.token)
}

func invalidOperandError(op *m.Token, found any, exp ...any) *RuntimeError {
	var expectedTypes []string
	for _, e := range exp {
		expectedTypes = append(expectedTypes, fmt.Sprintf("%T", e))
	}
	return  &RuntimeError{
			token: *op,
			message: fmt.Sprintf("operand(s) must be of type %s, but found a %T", strings.Join(expectedTypes, ","), found),
		}
}

func stringify(v any) string {
	return fmt.Sprintf("%v", v)
}

func NewInterp(errReporter lerr.Reporter) *Interp {
	return &Interp{
		errReporter: errReporter,
	}
}
