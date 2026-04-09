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
		w := " at end"
		r.report(t.Line, &w, message)
	} else {
		w := (" at '" + t.Lexeme + "'")
		r.report(t.Line, &w, message)
	}
}

func (r *Reporter) RuntimeError(t *models.Token, message *string) {
	fmt.Printf("%s %c[line \"%d\"]\n", *message, '\n', t.Line)
	r.HadRuntimeErr = true
}

func (r *Reporter) report(line int, where, message *string) {
	fmt.Println(
		"[line ", line, "] error", *where, ": ", *message,
	)
	r.HadErr = true
}

func NewReporter(initErr bool) *Reporter {
	return &Reporter{
		HadErr: initErr,
		HadRuntimeErr: initErr,
	}
}