// pretty printer

package ast

import "strings"

type ASTPrinter struct{}

func (p *ASTPrinter) Print(stmt Stmt) (string, error) {
	return AcceptStmt(stmt, p)
}

// =================================================================================================
// Statement Methods
// =================================================================================================

func (p *ASTPrinter) VisitExprStmtNodeStmt(stmt *ExprStmtNode) (string, error) {
  return AcceptExpr(stmt.Expr, p)
}

func (p *ASTPrinter) VisitPrintNodeStmt(stmt *PrintNode) (string, error) {
  return p.parenthesize("print", stmt.Expr), nil
}

func (p *ASTPrinter) VisitDeclNodeStmt(stmt *DeclNode) (string, error) {
  if stmt.Initializer != nil {
    val, _ := AcceptExpr(stmt.Ident, p)
    return p.parenthesize(val, stmt.Initializer), nil
  } else {
    return stmt.Ident.Name.Lexeme, nil
  }
}

func (p *ASTPrinter) VisitBlockNodeStmt(stmt *BlockNode) (string, error) {
  var ret []string
  for _, stmt := range stmt.Stmts {
    val, _ := AcceptStmt(stmt, p)
    ret = append(ret, val)
  }
  return strings.Join(ret, ",\n"), nil
}

// =================================================================================================
// Expression methods
// =================================================================================================

func (p *ASTPrinter) VisitBinaryNodeExpr(expr *BinaryNode) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, expr.Left, expr.Right), nil
}

func (p *ASTPrinter) VisitGroupingNodeExpr(expr *GroupingNode) (string, error) {
	return p.parenthesize("group", expr.Expr), nil
}

func (p *ASTPrinter) VisitLiteralNodeExpr(expr *LiteralNode) (string, error) {
	if (expr.Value == nil) { return "nil", nil }
	return expr.Value.Lexeme, nil
}

func (p *ASTPrinter) VisitUnaryNodeExpr(expr *UnaryNode) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, expr.Right), nil;
}

func (p *ASTPrinter) VisitIdentNodeExpr(expr *IdentNode) (string, error) {
  return expr.Name.Lexeme, nil
}

func (p *ASTPrinter) VisitAssignNodeExpr(expr *AssignNode) (string, error) {
  val, _ := AcceptExpr(expr.Ident, p)
  return p.parenthesize(val, expr.Value), nil
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
    // StringBuilder builder = new StringBuilder();
		builder := ""

    builder += "(" + name;
    for _, expr := range exprs {
      builder += " "
      new, _ := AcceptExpr(expr, p)
      builder += new
    }
    builder += ")"

    return builder;
  }

func NewASTPringer() *ASTPrinter {
	return &ASTPrinter{}
}