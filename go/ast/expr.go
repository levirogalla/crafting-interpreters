// generate file, do not edit

package ast

import (
	"crafting-interpreters/models"
	"fmt"
)

type Expr interface {
	isExpr()
}	
type ExprVisitor[T any] interface {
  VisitBinaryNodeExpr(expr *BinaryNode) (T, error)
  VisitGroupingNodeExpr(expr *GroupingNode) (T, error)
  VisitLiteralNodeExpr(expr *LiteralNode) (T, error)
  VisitUnaryNodeExpr(expr *UnaryNode) (T, error)
  VisitIdentNodeExpr(expr *IdentNode) (T, error)
}
type BinaryNode struct {
  Left Expr
  Op *models.Token
  Right Expr
}
func (BinaryNode) isExpr() {}

type GroupingNode struct {
  Expr Expr
}
func (GroupingNode) isExpr() {}

type LiteralNode struct {
  Value *models.Token
}
func (LiteralNode) isExpr() {}

type UnaryNode struct {
  Op *models.Token
  Right Expr
}
func (UnaryNode) isExpr() {}

type IdentNode struct {
  Name *models.Token
}
func (IdentNode) isExpr() {}

func AcceptExpr[T any](expr Expr, visitor ExprVisitor[T]) (T, error) {
  switch e := expr.(type) {
  case BinaryNode: return visitor.VisitBinaryNodeExpr(&e)
  case *BinaryNode: return visitor.VisitBinaryNodeExpr(e)
  case GroupingNode: return visitor.VisitGroupingNodeExpr(&e)
  case *GroupingNode: return visitor.VisitGroupingNodeExpr(e)
  case LiteralNode: return visitor.VisitLiteralNodeExpr(&e)
  case *LiteralNode: return visitor.VisitLiteralNodeExpr(e)
  case UnaryNode: return visitor.VisitUnaryNodeExpr(&e)
  case *UnaryNode: return visitor.VisitUnaryNodeExpr(e)
  case IdentNode: return visitor.VisitIdentNodeExpr(&e)
  case *IdentNode: return visitor.VisitIdentNodeExpr(e)
  default: panic(fmt.Sprintf("visitor not implemented for %T", expr))
  }
}
	