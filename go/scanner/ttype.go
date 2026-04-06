package scanner

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