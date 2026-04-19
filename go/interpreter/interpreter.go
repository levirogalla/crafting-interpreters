package interpreter

import (
	"crafting-interpreters/ast"
	"crafting-interpreters/environ"
	lerr "crafting-interpreters/error"
	m "crafting-interpreters/models"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type Interp struct {
	errReporter lerr.Reporter
	environ *environ.Environ
}

func (i *Interp) Interpret(stmts []ast.Stmt) {
	var err error
	for _, stmt := range stmts {
		err = i.execute(stmt)
		if err != nil {
			goto handleError
		}
	}
	return

handleError:
	t := (err.(lerr.RuntimeError).Token)
	m := (err.(lerr.RuntimeError).Message)
	i.errReporter.RuntimeError(&t, &m)
}

func (i *Interp) evaluate(expr ast.Expr) (m.Ltype, error) {
	return ast.AcceptExpr[m.Ltype](expr, i)
}

func (i *Interp) execute(stmt ast.Stmt) error {
	_, err := ast.AcceptStmt(stmt, i)
	return err
}

func (i *Interp) execBlock(stmts []ast.Stmt, environ *environ.Environ) error {
	prev := i.environ;
	defer func() { i.environ = prev }()

	i.environ = environ
	for _, stmt := range stmts {
		err := i.execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// =================================================================================================
// Statement methods
// =================================================================================================

func (i *Interp) VisitExprStmtNodeStmt(stmt *ast.ExprStmtNode) (int, error) {
	_, err := i.evaluate(stmt.Expr)
	return 0, err
}

func (i *Interp) VisitPrintNodeStmt(stmt *ast.PrintNode) (int, error) {
	val, err := i.evaluate(stmt.Expr)
	fmt.Println(stringify(val))
	return 0, err
}

func (i *Interp) VisitDeclNodeStmt(stmt *ast.DeclNode) (int, error) {
	var value m.Ltype
	var err error
	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
	}
	i.environ.Define(stmt.Ident.Name.Lexeme, value)
	return 0, err
}

func (i *Interp) VisitBlockNodeStmt(stmt *ast.BlockNode) (int, error) {
	err := i.execBlock(stmt.Stmts, environ.NewEnviron(i.environ));
	return 0, err
}

func (i *Interp) VisitIfNodeStmt(stmt *ast.IfNode) (int, error) {
	if isTruthy(stmt.Cond) {
		return 0,i.execute(stmt.Then)
	} else {
		return 0,i.execute(stmt.Else)
	}
}

// =================================================================================================
// Expression methods
// =================================================================================================

func (i *Interp) VisitLogicNodeExpr(expr *ast.LogicNode) (m.Ltype, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Op.TType == m.Or {
		if isTruthy(left) { return left, nil } 	
	} else {
		if !isTruthy(left) { return left, nil }
	}

	return i.evaluate(expr.Right)
}

func (i *Interp) VisitIdentNodeExpr(expr *ast.IdentNode) (m.Ltype, error) {
	return i.environ.Get(*expr.Name)
}

func (*Interp) VisitLiteralNodeExpr(expr *ast.LiteralNode) (m.Ltype, error) {
	return expr.Value.Lit, nil
}

func (i *Interp) VisitGroupingNodeExpr(expr *ast.GroupingNode) (m.Ltype, error) {
	return i.evaluate(expr.Expr)
}

func (i *Interp) VisitBinaryNodeExpr(expr *ast.BinaryNode) (m.Ltype, error) {
	var zNum m.Lnum
	var zStr m.Lstring
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Op.TType {
	case m.Minus:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
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
		if err != nil {
			return nil, err
		}
		return l / r, nil
	case m.Star:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
		return l * r, nil
	case m.GT:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
		return m.Lbool(l > r), nil
	case m.GTE:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
		return m.Lbool(l >= r), nil
	case m.LT:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
		return m.Lbool(l < r), nil
	case m.LTE:
		l, r, err := checkTypes[m.Lnum](expr.Op, left, right)
		if err != nil {
			return nil, err
		}
		return m.Lbool(l <= r), nil
	case m.Eq:
		return isEq(left, right), nil
	case m.Neq:
		return !isEq(left, right), nil
	default:
		panic("unreachable")
	}
}

func (i *Interp) VisitUnaryNodeExpr(expr *ast.UnaryNode) (m.Ltype, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Op.TType {
	case m.Bang:
		return !isTruthy(right), nil
	case m.Minus:
		r, err := checkType[m.Lnum](expr.Op, right)
		if err != nil {
			return nil, err
		}
		return -r, nil
	default:
		panic("unreachable: switch unary operator type")
	}
}

func (i *Interp) VisitAssignNodeExpr(expr *ast.AssignNode) (m.Ltype, error) {
	val, err := i.evaluate(expr.Value);
	if err != nil {
		return nil, err
	}
	err = i.environ.Assign(*expr.Ident.Name, val)
	return val, err
}

// =================================================================================================
// Utils
// =================================================================================================

func isTruthy(expr any) m.Lbool {
	if expr == nil {
		return m.Lbool(false)
	}
	switch e := expr.(type) {
	case bool:
		return m.Lbool(e)
	default:
		return m.Lbool(true)
	}
}

func isEq(l any, r any) m.Lbool {
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
			return m.Lbool(slices.Equal(l_, r_))
		}
		return false
	default:
		return m.Lbool(reflect.DeepEqual(l, r))
	}
}

func checkTypes[T m.Ltype](op *m.Token, a, b m.Ltype) (a_ T, b_ T, err error) {
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

func checkType[T m.Ltype](op *m.Token, a m.Ltype) (a_ T, err error) {
	var ok bool
	var zero T
	a_, ok = a.(T)
	if !ok {
		return zero, invalidOperandError(op, a, zero)
	}
	return a_, nil
}



func invalidOperandError(op *m.Token, found m.Ltype, exp ...m.Ltype) lerr.RuntimeError {
	var expectedTypes []string
	for _, e := range exp {
		expectedTypes = append(expectedTypes, m.StringifyLType(e))
	}
	var found_ string
	if found == nil {
		found_ = "nil"
	} else {
		found_ = fmt.Sprintf("'%v' of type '%s'", found, m.StringifyLType(found))
	}

	return lerr.RuntimeError{
		Token:   *op,
		Message: fmt.Sprintf("operand(s) must be of type '%s', but found %s.", strings.Join(expectedTypes, ", "), found_),
	}
}

func stringify(v any) string {
	return fmt.Sprintf("%v", v)
}

func NewInterp(errReporter lerr.Reporter) *Interp {
	return &Interp{
		errReporter: errReporter,
		environ: environ.NewEnviron(nil),
	}
}
