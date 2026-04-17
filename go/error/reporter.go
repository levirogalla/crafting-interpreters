package error

import (
	"crafting-interpreters/models"
	"fmt"
)

type Reporter struct {
	HadErr bool
	HadRuntimeErr bool
}

func (r *Reporter) ScanError(line int, message *string) {
	w := ""
	r.report(line, &w, message)
}

func (r *Reporter) ParseError(t *models.Token, message *string) {
	if t.TType == models.EOF {
		w := "at end"
		r.report(t.Line, &w, message)
	} else {
		w := ("at '" + t.Lexeme + "'")
		r.report(t.Line, &w, message)
	}
}

func (r *Reporter) RuntimeError(t *models.Token, message *string) {
	fmt.Printf("[line %d]: %s %c", t.Line, *message, '\n')
	r.HadRuntimeErr = true
}

func (r *Reporter) report(line int, where, message *string) {
	fmt.Printf(
		"[line %d] error %s: %s\n", line, *where, *message,
	)
	r.HadErr = true
}

func NewReporter(initErr bool) *Reporter {
	return &Reporter{
		HadErr: initErr,
		HadRuntimeErr: initErr,
	}
}

type RuntimeError struct {
	Token   models.Token
	Message string
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", r.Message, r.Token)
}

type ParseError struct {}
func (pErr ParseError) Error() string {
	return "parse error"
}