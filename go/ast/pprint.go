// pretty printer

package ast

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr Expr) string {
	return Accept(expr, p)
}

func (p *ASTPrinter) visitBinaryExpr(expr *Binary) string {
	return p.parenthesize(expr.Op.Lexeme, expr.Left, expr.Right)
}

func (p *ASTPrinter) visitGroupingExpr(expr *Grouping) string {
	return p.parenthesize("group", expr.Expr);
}

func (p *ASTPrinter) visitLiteralExpr(expr *Literal) string {
	if (expr.Value == nil) { return "nil" }
	return expr.Value.Lexeme
}

func (p *ASTPrinter) visitUnaryExpr(expr *Unary) string {
	return p.parenthesize(expr.Op.Lexeme, expr.Right);
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
    // StringBuilder builder = new StringBuilder();
		builder := ""

    builder += "(" + name;
    for _, expr := range exprs {
      builder += " "
      builder += Accept(expr, p)
    }
    builder += ")"

    return builder;
  }

func NewASTPringer() *ASTPrinter {
	return &ASTPrinter{}
}