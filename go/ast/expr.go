// generate file, do not edit

package ast

import (
	"crafting-interpreters/scanner"
)

type Expr interface {
	isExpr()
}	
type Visitor[T any] interface {
  visitBinaryExpr(expr *Binary) T
  visitGroupingExpr(expr *Grouping) T
  visitLiteralExpr(expr *Literal) T
  visitUnaryExpr(expr *Unary) T
}
type Binary struct {
  Left Expr
  Op *scanner.Token
  Right Expr
}
func (Binary) isExpr() {}

type Grouping struct {
  Expr Expr
}
func (Grouping) isExpr() {}

type Literal struct {
  Value *scanner.Token
}
func (Literal) isExpr() {}

type Unary struct {
  Op *scanner.Token
  Right Expr
}
func (Unary) isExpr() {}

func Accept[T any](expr Expr, visitor Visitor[T]) T {
  switch e := expr.(type) {
  case Binary: return visitor.visitBinaryExpr(&e)
  case Grouping: return visitor.visitGroupingExpr(&e)
  case Literal: return visitor.visitLiteralExpr(&e)
  case Unary: return visitor.visitUnaryExpr(&e)
  default: panic("unreachable")
  }
}
	