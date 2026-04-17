//go:generate ./cmd/generate-ast/generate-ast ./ast/expr.grammar -o ./ast/expr.go
//go:generate ./cmd/generate-ast/generate-ast ./ast/stmt.grammar -o ./ast/stmt.go
package main

import (
	"bufio"
	"crafting-interpreters/ast"
	"crafting-interpreters/error"
	"crafting-interpreters/interpreter"
	"crafting-interpreters/parser"
	"crafting-interpreters/scanner"
	"fmt"
	"os"
)

const testInputFile = "/Users/levirogalla/Projects/lib/crafting-interpreters/main.lox"

var errReporter = error.NewReporter(false)

func main() {
	args := os.Args
	// scanner.Main()
	// return
	// ASTPrinter := ast.NewASTPringer()
	// testTree := testAST()

	if len(args) > 2 {
		fmt.Println("usage: glox [script]")
	} else if len(args) == 2 {
		if args[1] == "test" {
			runFile(testInputFile)
		} else {
			runFile(args[1])
		}
	} else {
		runPrompt()
	}
}

func runFile(fp string) {
	bs, err := os.ReadFile(fp)
	if err != nil {
		fmt.Println(err)
	}
	run(string(bs))
	if errReporter.HadErr { os.Exit(65) }
	if errReporter.HadRuntimeErr { os.Exit(70) }
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			run(line)
			errReporter.HadErr = false
			errReporter.HadRuntimeErr = false
		}
		fmt.Print("> ")
	}
}

var interp = interpreter.NewInterp(*errReporter)
var printer = ast.NewASTPringer()
func run(i string) {
	scanner := scanner.NewScanner(i, *errReporter)
	ts := scanner.ScanTokens()
	parser := parser.NewParser(ts)
	stmts, err := parser.Parse()
	for _, stmt := range stmts {
		str, _ := printer.Print(stmt)
		fmt.Println(str)
		_ = str
	}
	_ = err
	interp.Interpret(stmts)
}
