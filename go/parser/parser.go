package parser

import (
	"crafting-interpreters/ast"
	lerr "crafting-interpreters/error"
	"crafting-interpreters/models"
	"fmt"
)

type Parser struct {
	tokens []models.Token
	current int
	errReporter lerr.Reporter
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	stmts := []ast.Stmt{}
	for (!p.done()) {
		stmt, err := p.declaration()
		if err != nil {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) declaration() (ast.Stmt, error) {
	var varDecl, stmt ast.Stmt
	var err error
	if p.match(models.Var) {
		varDecl, err = p.varDeclaration()
		if err != nil { 
			goto sync
		}
		return varDecl, nil
	}
	stmt, err = p.statement()
	if err != nil { 
		goto sync
	}
	return stmt, nil
sync:
	p.sync()
	return nil, err
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(models.Ident, "expected variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(models.Assign) {
		i, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = i
	}
	p.consume(models.Semicol, "expected ';' after variable declaration")
	return ast.DeclNode{
		Ident: ast.IdentNode{
			Name: name,
		},
		Initializer: initializer,
	}, nil
}


func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(models.Print) { 
		return p.printStatement()
	}
	return p.exprStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil { return nil, err }
	p.consume(models.Semicol, "expected ';' after value for print.")
	return &ast.PrintNode{
		Expr: value,
	}, nil
}

func (p *Parser) exprStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil { return nil, err }
	p.consume(models.Semicol, "expected ';' after value for expression.")
	return &ast.ExprStmtNode{
		Expr: value,
	}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignent()
}

func (p *Parser) assignent() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(models.Assign) {
		eq := p.previous();
		rval, _ := p.assignent();

		switch e := expr.(type) {
		case *ast.IdentNode:
			return &ast.AssignNode{
				Ident: e,
				Value: rval,
			}, nil
		}

		p.error(eq, fmt.Sprintf("invalid assignment target. %T", expr))
	}
	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil { return nil, err }

	for (p.match(models.Neq, models.Eq)) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil { return nil, err }
		expr = &ast.BinaryNode{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil { return nil, err }

	for p.match(models.GT, models.GTE, models.LT, models.LTE) {
		op := p.previous()
		right, err := p.term()
		if err != nil { return nil, err }
		expr = &ast.BinaryNode{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil { return nil, err }

	for p.match(models.Minus, models.Plus) {
		op := p.previous()
		right, err := p.factor()
		if err != nil { return nil, err }
		expr = &ast.BinaryNode{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil { return nil, err }

	for p.match(models.Slash, models.Star) {
		op := p.previous()
		right, err := p.unary()
		if err != nil { return nil, err }
		expr = &ast.BinaryNode{Left: expr, Op: op, Right: right}
	}

	return expr, nil

}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(models.Bang, models.Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil { return nil, err }	
		return &ast.UnaryNode{Op: op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	tk := p.advance()
	switch tk.TType {
	case models.False, models.True, models.Nil, models.Num, models.Str:
		return &ast.LiteralNode{Value: tk}, nil

	case models.LParen:
		expr, err := p.expression()
		if err != nil { return nil, err }	
		p.consume(models.RParen, "expect ')' after expression")
		return &ast.GroupingNode{Expr: expr}, nil

	case models.Ident:
		return &ast.IdentNode{
			Name: tk,
		}, nil
	}

	return nil, p.error(tk, fmt.Sprintf("expected expression, found '%s'.", tk.Lexeme))
}

func (p *Parser) sync() {
	p.advance()
	
	for !p.done() {
		if p.previous().TType == models.Semicol { return }

		switch p.peek().TType {
		case models.Class, models.Fn, models.Var, models.For, models.If, models.While, models.Print, models.Return:
			return
		}

		p.advance()
	}
}

// =================================================================================================
// utils
// =================================================================================================

func (p *Parser) match(ttypes ...models.TokenType) bool {
	for _, t := range ttypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tt models.TokenType) bool {
	if p.done() { return false }
	return p.peek().TType == tt
}

func (p *Parser) advance() *models.Token {
	if !p.done() { p.current++ }
	return p.previous()
}

func (p *Parser) done() bool {
	return p.peek().TType == models.EOF
}

func (p *Parser) peek() models.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *models.Token {
	return &p.tokens[p.current - 1]
}

func (p *Parser) consume(tt models.TokenType, message string) (*models.Token, error) {
	if p.check(tt) { return p.advance(), nil }
	t := p.peek()
	return nil, p.error(&t, message)
}

func (p *Parser) error(t *models.Token, message string) *lerr.ParseError {
	p.errReporter.ParseError(t, &message)
	return &lerr.ParseError{}
}

func NewParser(ts *[]models.Token) *Parser {
	return &Parser{
		tokens: *ts,
		current: 0,
	}
}

