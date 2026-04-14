package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func main() {
	args := os.Args

	outputDir, err := os.Getwd()
	e(err)

	// parse args

	var input string
	var super string
	var output string
	switch len(args) {
	case 2:
		input = args[1]
		super = strings.TrimSuffix(filepath.Base(input), ".grammar")
		output = filepath.Join(outputDir, filepath.Base(input))
		output = strings.TrimSuffix(output, ".grammar") + ".go"
	case 4:
		input = args[1]
		super = strings.TrimSuffix(filepath.Base(input), ".grammar")
		output = args[3]
	default:
		incorrectUsage()
	}
	super = capitalize(super)

	// reset file

	err = os.WriteFile(output, []byte{}, 0644)
	e(err)

	// preamble

	f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e(err)
	defer f.Close()

	fmt.Fprintf(f,
		`// generate file, do not edit

package ast

import (
	"crafting-interpreters/models"
	"fmt"
)

type %s interface {
	is%s()
}	
`, super, super)
	grammar, err := os.ReadFile(input)
	e(err)
	lines := strings.Split(string(grammar), "\n")

	// visitor pattern

	fmt.Fprintf(f, "type %sVisitor[T any] interface {\n", super)
	for i, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) != 2 {
			exit(fmt.Sprintf("error parsing on line %d", i))
		}

		name := strings.TrimSpace(parts[0])

		fmt.Fprintf(f, "  Visit%s%s(expr *%s) (T, error)\n", name, super, name)
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

		fmt.Fprintf(f, "}\nfunc (%s) is%s() {}\n", name, super)
// 		fmt.Fprintf(f, `func (self %s[T]) accept(visitor Visitor[T]) T {
// 	return visitor.visit%s%s(self)
// }
// `, name, name, "Expr")
		fmt.Fprintln(f)
	}

	fmt.Fprintf(f, "func Accept%s[T any](expr %s, visitor %sVisitor[T]) (T, error) {\n  switch e := expr.(type) {\n", super, super, super)
	for i, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) != 2 {
			exit(fmt.Sprintf("error parsing on line %d", i))
		}

		name := strings.TrimSpace(parts[0])
		fmt.Fprintf(f, "  case %s: return visitor.Visit%s%s(&e)\n", name, name, super)
		fmt.Fprintf(f, "  case *%s: return visitor.Visit%s%s(e)\n", name, name, super)
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

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}