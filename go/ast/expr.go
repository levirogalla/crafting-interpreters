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
  VisitAssignNodeExpr(expr *AssignNode) (T, error)
  VisitBinaryNodeExpr(expr *BinaryNode) (T, error)
  VisitGroupingNodeExpr(expr *GroupingNode) (T, error)
  VisitLiteralNodeExpr(expr *LiteralNode) (T, error)
  VisitUnaryNodeExpr(expr *UnaryNode) (T, error)
  VisitIdentNodeExpr(expr *IdentNode) (T, error)
  VisitLogicNodeExpr(expr *LogicNode) (T, error)
}
type AssignNode struct {
  Ident *IdentNode
  Value Expr
}
func (AssignNode) isExpr() {}

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

type LogicNode struct {
  Left Expr
  Op *models.Token
  Right Expr
}
func (LogicNode) isExpr() {}

func AcceptExpr[T any](expr Expr, visitor ExprVisitor[T]) (T, error) {
  switch e := expr.(type) {
  case AssignNode: return visitor.VisitAssignNodeExpr(&e)
  case *AssignNode: return visitor.VisitAssignNodeExpr(e)
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
  case LogicNode: return visitor.VisitLogicNodeExpr(&e)
  case *LogicNode: return visitor.VisitLogicNodeExpr(e)
  default: panic(fmt.Sprintf("visitor not implemented for %T", expr))
  }
}
	