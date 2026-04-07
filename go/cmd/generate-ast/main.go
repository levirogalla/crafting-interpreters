package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args

	outputDir, err := os.Getwd()
	e(err)

	// parse args

	var input string
	var output string
	switch len(args) {
	case 2:
		input = args[1]
		output = filepath.Join(outputDir, filepath.Base(input))
		output = strings.TrimSuffix(output, ".grammar") + ".go"
	case 4:
		input = args[1]
		output = args[3]
	default:
		incorrectUsage()
	}

	// reset file

	err = os.WriteFile(output, []byte{}, 0644)
	e(err)

	// preamble

	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e(err)
	defer f.Close()

	f.WriteString(
		`// generate file, do not edit

package ast

import (
	"crafting-interpreters/models"
	"fmt"
)

type Expr interface {
	isExpr()
}	
`)
	grammar, err := os.ReadFile(input)
	e(err)
	lines := strings.Split(string(grammar), "\n")

	// visitor pattern

	fmt.Fprintf(f, "type Visitor[T any] interface {\n")
	for i, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) != 2 {
			exit(fmt.Sprintf("error parsing on line %d", i))
		}

		name := strings.TrimSpace(parts[0])

		fmt.Fprintf(f, "  visit%s%s(expr *%s) T\n", name, "Expr", name)
	}
	fmt.Fprintf(f, "}\n")

	// ast node types

	for i, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) != 2 {
			exit(fmt.Sprintf("error parsing on line %d", i))
		}

		name := strings.TrimSpace(parts[0])
		fmt.Fprintf(f, "type %s struct {\n", name)

		right := strings.TrimSpace(parts[1])

		fields := strings.Split(right, ",")

		for j, field := range fields {
			field = strings.TrimSpace(field)

			parts := strings.Fields(field)
			if len(parts) != 2 {
				exit(fmt.Sprintf("error parsing on line %d-%d", i, j))
			}

			typ := parts[0]
			varName := parts[1]

			fmt.Fprintf(f, "  %s %s\n", varName, typ)
		}

		fmt.Fprintf(f, "}\nfunc (%s) isExpr() {}\n", name)
// 		fmt.Fprintf(f, `func (self %s[T]) accept(visitor Visitor[T]) T {
// 	return visitor.visit%s%s(self)
// }
// `, name, name, "Expr")
		fmt.Fprintln(f)
	}

	fmt.Fprintf(f, "func Accept[T any](expr %s, visitor Visitor[T]) T {\n  switch e := expr.(type) {\n", "Expr")
	for i, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) != 2 {
			exit(fmt.Sprintf("error parsing on line %d", i))
		}

		name := strings.TrimSpace(parts[0])
		fmt.Fprintf(f, "  case %s: return visitor.visit%s%s(&e)\n", name, name, "Expr")
		fmt.Fprintf(f, "  case *%s: return visitor.visit%s%s(e)\n", name, name, "Expr")
	}
	fmt.Fprint(f, `  default: panic(fmt.Sprintf("visitor not implemented for %T", expr))
  }
}
	`)
}

func incorrectUsage() {
	fmt.Println("usage: generate-ast <input> [-o]")
	os.Exit(65)
}

func e(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(65)
	}
}

func exit(err string) {
	fmt.Println(err)
	os.Exit(65)
}
