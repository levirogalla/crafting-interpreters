// generate file, do not edit

package ast

import (
	"crafting-interpreters/models"
	"fmt"
)

type Expr interface {
	isExpr()
}	
type Visitor[T any] interface {
  VisitBinaryExpr(expr *Binary) (T, error)
  VisitGroupingExpr(expr *Grouping) (T, error)
  VisitLiteralExpr(expr *Literal) (T, error)
  VisitUnaryExpr(expr *Unary) (T, error)
}
type Binary struct {
  Left Expr
  Op *models.Token
  Right Expr
}
func (Binary) isExpr() {}

type Grouping struct {
  Expr Expr
}
func (Grouping) isExpr() {}

type Literal struct {
  Value *models.Token
}
func (Literal) isExpr() {}

type Unary struct {
  Op *models.Token
  Right Expr
}
func (Unary) isExpr() {}

func Accept[T any](expr Expr, visitor Visitor[T]) (T, error) {
  switch e := expr.(type) {
  case Binary: return visitor.VisitBinaryExpr(&e)
  case *Binary: return visitor.VisitBinaryExpr(e)
  case Grouping: return visitor.VisitGroupingExpr(&e)
  case *Grouping: return visitor.VisitGroupingExpr(e)
  case Literal: return visitor.VisitLiteralExpr(&e)
  case *Literal: return visitor.VisitLiteralExpr(e)
  case Unary: return visitor.VisitUnaryExpr(&e)
  case *Unary: return visitor.VisitUnaryExpr(e)
  default: panic(fmt.Sprintf("visitor not implemented for %T", expr))
  }
}
	