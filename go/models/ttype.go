package models

import "fmt"

type TokenType int


func (tt TokenType) String() string {
	names := [...]string{
		"Inval",
		"LParen",
		"RParen",
		"LBrace",
		"RBrace",
		"Comma",
		"Dot",
		"Minus",
		"Plus",
		"Semicol",
		"Slash",
		"Star",

		"Bang",
		"Neq",
		"Asign",
		"Eq",
		"GT",
		"LT",
		"GTE",
		"LTE",

		"Ident",
		"Str",
		"Num",

		"And",
		"Class",
		"Else",
		"False",
		"Fn",
		"For",
		"If",
		"Nil",
		"Or",
		"Print",
		"Return",
		"Super",
		"Self",
		"True",
		"Var",
		"While",

		"EOF",
	}

	if int(tt) < 0 || int(tt) >= len(names) {
		return "Unknown"
	}
	return names[tt]
}
const (
	// single char tokens
	Inval TokenType = iota
	LParen 
	RParen
	LBrace
	RBrace
	Comma
	Dot
	Minus
	Plus
	Semicol
	Slash
	Star

	// one or two char tokens
	Bang
	Neq
	Asign
	Eq
	GT
	LT
	GTE
	LTE

	// literals
	Ident
	Str
	Num

	// kewwords
	And
	Class
	Else
	False
	Fn
	For
	If
	Nil
	Or
	Print
	Return
	Super
	Self
	True
	Var
	While

	// special
	EOF
)

type Token struct {
	TType  TokenType
	Lexeme string
	Lit    Ltype 
	Line   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token { type: %s, lexeme: %s, lit: %T{%v}, line: %d }",
		t.TType, t.Lexeme, t.Lit, t.Lit, t.Line)
}


func NewToken(ttype TokenType, lexeme string, lit Ltype, line int) *Token {
	return &Token{TType: ttype, Lexeme: lexeme, Lit: lit, Line: line}
}

type Lnum float64
func (Lnum) isLType() {}

type Lstring string
func (Lstring) isLType() {}

type Lbool bool
func (Lbool) isLType() {}

type Ltype interface {
	isLType()
}

func StringifyLType(ltype Ltype) string {
	switch ltype.(type) {
	case Lnum, *Lnum:
		return "num"
	case Lstring, *Lstring:
		return "string"
	case Lbool, *Lbool:
		return "bool"
	case nil:
		return "nil"
	default: panic("unsupporeted ltype")
	}
}