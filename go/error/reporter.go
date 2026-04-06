package error

import "fmt"

type Reporter struct {
	hadErr bool
}

func (r Reporter) Error(line int, message string) {
	r.Report(line, "", message)
}

func (r Reporter) Report(line int, where, message string) {
	fmt.Println(
		"[line ", line, "] error", where, ": ", message,
	)
	r.hadErr = true
}

func (r Reporter) HadError() bool {
	return r.hadErr
}

func NewReporter(initErr bool) *Reporter {
	return &Reporter{
		hadErr: initErr,
	}
}