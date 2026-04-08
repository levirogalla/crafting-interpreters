// pretty printer

package ast

type ASTPrinter struct{}

func (p *ASTPrinter) Print(expr Expr) (string, error) {
	return Accept(expr, p)
}

func (p *ASTPrinter) VisitBinaryExpr(expr *Binary) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, expr.Left, expr.Right), nil
}

func (p *ASTPrinter) VisitGroupingExpr(expr *Grouping) (string, error) {
	return p.parenthesize("group", expr.Expr), nil
}

func (p *ASTPrinter) VisitLiteralExpr(expr *Literal) (string, error) {
	if (expr.Value == nil) { return "nil", nil }
	return expr.Value.Lexeme, nil
}

func (p *ASTPrinter) VisitUnaryExpr(expr *Unary) (string, error) {
	return p.parenthesize(expr.Op.Lexeme, expr.Right), nil;
}

func (p *ASTPrinter) parenthesize(name string, exprs ...Expr) string {
    // StringBuilder builder = new StringBuilder();
		builder := ""

    builder += "(" + name;
    for _, expr := range exprs {
      builder += " "
      new, _ := Accept(expr, p)
      builder += new
    }
    builder += ")"

    return builder;
  }

func NewASTPringer() *ASTPrinter {
	return &ASTPrinter{}
}