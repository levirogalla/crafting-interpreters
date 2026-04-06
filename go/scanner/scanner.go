package scanner

import (
	"crafting-interpreters/error"
	"fmt"
	"strconv"
)

type Token struct {
	ttype  TokenType
	Lexeme string
	lit    any 
	line   int
}

func (t Token) String() string {
	return fmt.Sprintf("Token { type: %s, lexeme: %s, lit: %T{%v}, line: %d }",
		t.ttype, t.Lexeme, t.lit, t.lit, t.line)
}


func NewToken(ttype TokenType, lexeme string, lit interface{}, line int) *Token {
	return &Token{ttype: ttype, Lexeme: lexeme, lit: lit, line: line}
}

type Scanner struct {
	start, current, line int
	ts                   []Token
	source               string
	errReporter          error.Reporter
}

func NewScanner(i string, r error.Reporter) *Scanner {
	return &Scanner{
		start:       0,
		current:     0,
		line:        0,
		ts:          []Token{},
		source:      i,
		errReporter: r,
	}
}

func (s *Scanner) done() bool {
	return s.current == len(s.source)
}

func (s *Scanner) ScanTokens() []Token {
	for !s.done() {
		s.start = s.current
		s.scanToken()
	}
	s.addSimpleToken(EOF)
	return s.ts
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addSimpleToken(LParen)
	case ')':
		s.addSimpleToken(RParen)
	case '{':
		s.addSimpleToken(LBrace)
	case '}':
		s.addSimpleToken(RBrace)
	case ',':
		s.addSimpleToken(Comma)
	case '.':
		s.addSimpleToken(Dot)
	case '-':
		s.addSimpleToken(Minus)
	case '+':
		s.addSimpleToken(Plus)
	case ';':
		s.addSimpleToken(Semicol)
	case '*':
		s.addSimpleToken(Star)

	case '!':
		if s.match('=') {
			s.addSimpleToken(Neq)
		} else {
			s.addSimpleToken(Bang)
		}

	case '=':
		if s.match('=') {
			s.addSimpleToken(Eq)
		} else {
			s.addSimpleToken(Asign)
		}

	case '<':
		if s.match('=') {
			s.addSimpleToken(LTE)
		} else {
			s.addSimpleToken(LT)
		}

	case '>':
		if s.match('=') {
			s.addSimpleToken(GTE)
		} else {
			s.addSimpleToken(GT)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.done() {
				s.advance()
			}
		} else {
			s.addSimpleToken(Slash)
		}

	case ' ':
	case '\r':
	case '\t':
		break

	case '\n':
		s.line++

	case '"':
		s.handleString()

	default:
		if isDigit(c) {
			s.handleNumber()
		} else if isAlpha(c) {
			s.handleIdent()
		} else {
			s.errReporter.Error(s.line, "unexpected char")
		}
	}
}

func (s *Scanner) advance() rune {
	c := s.peek()
	s.current++
	return c
}

func (s *Scanner) peek() rune {
	if s.done() { return 0 }
	c := rune(s.source[s.current])
	return c
}

func (s *Scanner) peek2() rune {
	if s.current+1 >= len(s.source) {
		return 0
	} else {
		return rune(s.source[s.current+1])
	}
}

func (s *Scanner) match(c rune) bool {
	if s.done() {
		return false
	}
	if s.peek() != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) addSimpleToken(t TokenType) {
	s.addToken(t, nil)
}

func (s *Scanner) addToken(t TokenType, lit interface{}) {
	text := s.source[s.start:s.current]
	s.ts = append(s.ts, *NewToken(t, text, lit, s.line))
}

// =================================================================================================
// more complicated stuff
// =================================================================================================

func (s *Scanner) handleString() {
	for !s.match('"') {
		if s.peek() != '\n' {
			s.line++
		}
	}

	val := s.source[s.start+1 : s.current-1]
	s.addToken(Str, val)
}

func (s *Scanner) handleNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peek2()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	val, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.errReporter.Error(s.line, "unable to parse number")
	}
	s.addToken(Num, val)
}

func (s *Scanner) handleIdent() {
	for isAlphaNum(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	switch text {
	case "and": s.addSimpleToken(And)
	case "class": s.addSimpleToken(Class)
	case "else": s.addSimpleToken(Else)
	case "false": s.addSimpleToken(False)
	case "for": s.addSimpleToken(For)
	case "fn": s.addSimpleToken(Fn)
	case "if": s.addSimpleToken(If)
	case "nil": s.addSimpleToken(Nil)
	case "or": s.addSimpleToken(Or)
	case "return": s.addSimpleToken(Return)
	case "super": s.addSimpleToken(Super)
	case "self": s.addSimpleToken(Self)
	case "true": s.addSimpleToken(True)
	case "var": s.addSimpleToken(Var)
	case "while": s.addSimpleToken(While)
	default:
		s.addToken(Ident, text)	
	}
}

// =================================================================================================
// Utils
// =================================================================================================

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

// func Main() {
// 	fmt.Println(isDigit('0'))
// 	fmt.Println(isDigit('9'))
// 	fmt.Println(isAlpha('a'))
// 	fmt.Println(isAlpha(' '))
// 	fmt.Println(isAlphaNum('0'))
// 	fmt.Println(isAlphaNum('9'))
// 	fmt.Println(isAlphaNum('a'))
// 	fmt.Println(isAlphaNum('Z'))
// }