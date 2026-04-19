// generate file, do not edit

package ast

import (
	"fmt"
)

type Stmt interface {
	isStmt()
}	
type StmtVisitor[T any] interface {
  VisitExprStmtNodeStmt(expr *ExprStmtNode) (T, error)
  VisitPrintNodeStmt(expr *PrintNode) (T, error)
  VisitDeclNodeStmt(expr *DeclNode) (T, error)
  VisitBlockNodeStmt(expr *BlockNode) (T, error)
  VisitIfNodeStmt(expr *IfNode) (T, error)
}
type ExprStmtNode struct {
  Expr Expr
}
func (ExprStmtNode) isStmt() {}

type PrintNode struct {
  Expr Expr
}
func (PrintNode) isStmt() {}

type DeclNode struct {
  Ident IdentNode
  Initializer Expr
}
func (DeclNode) isStmt() {}

type BlockNode struct {
  Stmts []Stmt
}
func (BlockNode) isStmt() {}

type IfNode struct {
  Cond Expr
  Then Stmt
  Else Stmt
}
func (IfNode) isStmt() {}

func AcceptStmt[T any](expr Stmt, visitor StmtVisitor[T]) (T, error) {
  switch e := expr.(type) {
  case ExprStmtNode: return visitor.VisitExprStmtNodeStmt(&e)
  case *ExprStmtNode: return visitor.VisitExprStmtNodeStmt(e)
  case PrintNode: return visitor.VisitPrintNodeStmt(&e)
  case *PrintNode: return visitor.VisitPrintNodeStmt(e)
  case DeclNode: return visitor.VisitDeclNodeStmt(&e)
  case *DeclNode: return visitor.VisitDeclNodeStmt(e)
  case BlockNode: return visitor.VisitBlockNodeStmt(&e)
  case *BlockNode: return visitor.VisitBlockNodeStmt(e)
  case IfNode: return visitor.VisitIfNodeStmt(&e)
  case *IfNode: return visitor.VisitIfNodeStmt(e)
  default: panic(fmt.Sprintf("visitor not implemented for %T", expr))
  }
}
	