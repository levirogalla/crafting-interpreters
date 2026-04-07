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

func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	if err != nil { return nil, err }	
	return expr, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil { return nil, err }

	for (p.match(models.Neq, models.Eq)) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil { return nil, err }
		expr = &ast.Binary{Left: expr, Op: &op, Right: right}
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
		expr = &ast.Binary{Left: expr, Op: &op, Right: right}
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
		expr = &ast.Binary{Left: expr, Op: &op, Right: right}
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
		expr = &ast.Binary{Left: expr, Op: &op, Right: right}
	}

	return expr, nil

}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(models.Bang, models.Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil { return nil, err }	
		return &ast.Unary{Op: &op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	tk := p.advance()
	switch tk.TType {
	case models.False, models.True, models.Nil, models.Num, models.Str:
		return &ast.Literal{Value: &tk}, nil

	case models.LParen:
		expr, err := p.expression()
		if err != nil { return nil, err }	
		p.consume(models.RParen, "expect ')' after expression")
		return &ast.Grouping{Expr: expr}, nil
	}

	// nxt := p.peek()
	return nil, p.error(&tk, fmt.Sprintf("expected expression %s", tk))
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
      // switch (peek().type) {
      //   case CLASS:
      //   case FUN:
      //   case VAR:
      //   case FOR:
      //   case IF:
      //   case WHILE:
      //   case PRINT:
      //   case RETURN:
      //     return;
      // }

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

func (p *Parser) advance() models.Token {
	if !p.done() { p.current++ }
	return p.previous()
}

func (p *Parser) done() bool {
	return p.peek().TType == models.EOF
}

func (p *Parser) peek() models.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() models.Token {
	return p.tokens[p.current - 1]
}

func (p *Parser) consume(tt models.TokenType, message string) (models.Token, error) {
	if p.check(tt) { return p.advance(), nil }
	t := p.peek()
	return models.Token{}, p.error(&t, message)
}

func (p *Parser) error(t *models.Token, message string) *ParseError {
	p.errReporter.ParseError(t, &message)
	return &ParseError{}
}

func NewParser(ts *[]models.Token) *Parser {
	return &Parser{
		tokens: *ts,
		current: 0,
	}
}

type ParseError struct {}
func (pErr ParseError) Error() string {
	return "parse error"
}

func bubble(err error) {}