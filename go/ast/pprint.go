// pretty printer

package ast

import "strings"

type ASTPrinter struct{}

func (p *ASTPrinter) PrintExpr(expr Expr) (string, error) {
  return AcceptExpr[string](expr, p)
}
func (p *ASTPrinter) PrintStmt(stmt Stmt) (string, error) {
  return AcceptStmt[string](stmt, p)
}

// =================================================================================================
// Statement Methods
// =================================================================================================

func (p *ASTPrinter) VisitExprStmtNodeStmt(stmt *ExprStmtNode) (string, error) {
  return AcceptExpr(stmt.Expr, p)
}

func (p *ASTPrinter) VisitPrintNodeStmt(stmt *PrintNode) (string, error) {
  return p.parenthesize("print", ExprStmtNode{stmt.Expr}), nil
}

func (p *ASTPrinter) VisitDeclNodeStmt(stmt *DeclNode) (string, error) {
  if stmt.Initializer != nil {
    val, _ := AcceptExpr(stmt.Ident, p)
    return p.parenthesize(val, ExprStmtNode{stmt.Initializer}), nil
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

func (p *ASTPrinter) VisitIfNodeStmt(stmt *IfNode) (string, error) {
  if stmt.Else == nil {
    return p.parenthesize("if", ExprStmtNode{stmt.Cond}, stmt.Then), nil
  } else {
    return p.parenthesize("if else", ExprStmtNode{stmt.Cond}, stmt.Then, stmt.Else), nil
  }
}

// =================================================================================================
// Expression methods
// =================================================================================================

func (p *ASTPrinter) VisitLogicNodeExpr(expr *LogicNode) (string, error) {
  return p.parenthesize(expr.Op.Lexeme, ExprStmtNode{expr.Left}, ExprStmtNode{expr.Right}), nil
}

func (p *ASTPrinter) VisitBinaryNodeExpr(expr *BinaryNode) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, ExprStmtNode{expr.Left}, ExprStmtNode{expr.Right}), nil
}

func (p *ASTPrinter) VisitGroupingNodeExpr(expr *GroupingNode) (string, error) {
	return p.parenthesize("group", ExprStmtNode{expr.Expr}), nil
}

func (p *ASTPrinter) VisitLiteralNodeExpr(expr *LiteralNode) (string, error) {
	if (expr.Value == nil) { return "nil", nil }
	return expr.Value.Lexeme, nil
}

func (p *ASTPrinter) VisitUnaryNodeExpr(expr *UnaryNode) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, ExprStmtNode{expr.Right}), nil;
}

func (p *ASTPrinter) VisitIdentNodeExpr(expr *IdentNode) (string, error) {
  return expr.Name.Lexeme, nil
}

func (p *ASTPrinter) VisitAssignNodeExpr(expr *AssignNode) (string, error) {
  val, _ := AcceptExpr(expr.Ident, p)
  return p.parenthesize(val, ExprStmtNode{expr.Value}), nil
}

func (p *ASTPrinter) parenthesize(name string, stmts ...Stmt) string {
    // StringBuilder builder = new StringBuilder();
		builder := ""

    builder += "(" + name;
    for _, expr := range stmts {
      builder += " "
      new, _ := AcceptStmt[string](expr, p)
      builder += new
    }
    builder += ")"

    return builder;
  }

func NewASTPringer() *ASTPrinter {
	return &ASTPrinter{}
}