package error

import (
	"crafting-interpreters/models"
	"fmt"
)

type Reporter struct {
	hadErr bool
}

func (r *Reporter) ScanError(line int, message *string) {
	w := ""
	r.report(line, &w, message)
}

func (r *Reporter) ParseError(t *models.Token, message *string) {
	if t.TType == models.EOF {
		w := " at end"
		r.report(t.Line, &w, message)
	} else {
		w := (" at '" + t.Lexeme + "'")
		r.report(t.Line, &w, message)
	}
}

func (r *Reporter) report(line int, where, message *string) {
	fmt.Println(
		"[line ", line, "] error", *where, ": ", *message,
	)
	r.hadErr = true
}

func (r *Reporter) HadError() bool {
	return r.hadErr
}

func NewReporter(initErr bool) *Reporter {
	return &Reporter{
		hadErr: initErr,
	}
}