package scanner

import (
	"crafting-interpreters/error"
	"crafting-interpreters/models"
	"strconv"
)



type Scanner struct {
	start, current, line int
	ts                   []models.Token
	source               string
	errReporter          error.Reporter
}

func NewScanner(i string, r error.Reporter) *Scanner {
	return &Scanner{
		start:       0,
		current:     0,
		line:        0,
		ts:          []models.Token{},
		source:      i,
		errReporter: r,
	}
}

func (s *Scanner) done() bool {
	return s.current == len(s.source)
}

func (s *Scanner) ScanTokens() *[]models.Token {
	for !s.done() {
		s.start = s.current
		s.scanToken()
	}
	s.addSimpleToken(models.EOF)
	return &s.ts
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addSimpleToken(models.LParen)
	case ')':
		s.addSimpleToken(models.RParen)
	case '{':
		s.addSimpleToken(models.LBrace)
	case '}':
		s.addSimpleToken(models.RBrace)
	case ',':
		s.addSimpleToken(models.Comma)
	case '.':
		s.addSimpleToken(models.Dot)
	case '-':
		s.addSimpleToken(models.Minus)
	case '+':
		s.addSimpleToken(models.Plus)
	case ';':
		s.addSimpleToken(models.Semicol)
	case '*':
		s.addSimpleToken(models.Star)

	case '!':
		if s.match('=') {
			s.addSimpleToken(models.Neq)
		} else {
			s.addSimpleToken(models.Bang)
		}

	case '=':
		if s.match('=') {
			s.addSimpleToken(models.Eq)
		} else {
			s.addSimpleToken(models.Assign)
		}

	case '<':
		if s.match('=') {
			s.addSimpleToken(models.LTE)
		} else {
			s.addSimpleToken(models.LT)
		}

	case '>':
		if s.match('=') {
			s.addSimpleToken(models.GTE)
		} else {
			s.addSimpleToken(models.GT)
		}

	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.done() {
				s.advance()
			}
		} else {
			s.addSimpleToken(models.Slash)
		}

	case ' ', '\r', '\t':
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
			m := "unexpected char"
			s.errReporter.ScanError(s.line, &m)
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

func (s *Scanner) addSimpleToken(t models.TokenType) {
	s.addToken(t, nil)
}

func (s *Scanner) addToken(t models.TokenType, lit models.Ltype) {
	text := s.source[s.start:s.current]
	s.ts = append(s.ts, *models.NewToken(t, text, lit, s.line))
}

// =================================================================================================
// more complicated stuff
// =================================================================================================

func (s *Scanner) handleString() {
	for !s.match('"') {
		if s.peek() != '\n' {
			s.line++
		}
		s.advance()
	}

	val := s.source[s.start+1 : s.current-1]
	s.addToken(models.Str, models.Lstring(val))
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
		m := "unable to parse number"
		s.errReporter.ScanError(s.line, &m)
	}
	s.addToken(models.Num, models.Lnum(val))
}

func (s *Scanner) handleIdent() {
	for isAlphaNum(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	switch text {
	case "and": s.addSimpleToken(models.And)
	case "class": s.addSimpleToken(models.Class)
	case "else": s.addSimpleToken(models.Else)
	case "false": s.addToken(models.False, models.Lbool(false))
	case "for": s.addSimpleToken(models.For)
	case "fn": s.addSimpleToken(models.Fn)
	case "if": s.addSimpleToken(models.If)
	case "nil": s.addToken(models.Nil, nil)
	case "or": s.addSimpleToken(models.Or)
	case "return": s.addSimpleToken(models.Return)
	case "super": s.addSimpleToken(models.Super)
	case "self": s.addSimpleToken(models.Self)
	case "true": s.addToken(models.True, models.Lbool(true))
	case "var": s.addSimpleToken(models.Var)
	case "while": s.addSimpleToken(models.While)
	case "print": s.addSimpleToken(models.Print)
	default:
		s.addToken(models.Ident, models.Lstring(text))	
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