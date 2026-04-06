//go:generate ./ast/generate-ast ./ast/expr.grammar -o ./ast/expr.go
package main

import (
	"bufio"
	"crafting-interpreters/ast"
	"crafting-interpreters/error"
	"crafting-interpreters/scanner"
	"fmt"
	"os"
)

const testInputFile = "/Users/levirogalla/Projects/lib/crafting-interpreters/main.lox"

var hadError = false
var errReporter = error.NewReporter(false)

func testAST() ast.Expr {
	expr := ast.Binary{
		Left: ast.Unary{
			Op: scanner.NewToken(scanner.Minus, "-", nil, 1),
			Right: ast.Literal{
				Value: scanner.NewToken(scanner.Num, "123", 123, 1),
			},
		},
		Op: scanner.NewToken(scanner.Star, "*", nil, 1),
		Right: ast.Grouping{
			Expr: ast.Literal{
				Value: scanner.NewToken(scanner.Num, "45.67", 45.67, 1),
			},
		},
	}
	return expr
	//   public static void main(String[] args) {
	//   Expr expression = new Expr.Binary(
	//       new Expr.Unary(
	//           new Token(TokenType.MINUS, "-", null, 1),
	//           new Expr.Literal(123)),
	//       new Token(TokenType.STAR, "*", null, 1),
	//       new Expr.Grouping(
	//           new Expr.Literal(45.67)));

	//   System.out.println(new AstPrinter().print(expression));
	// }
}

func main() {
	args := os.Args
	// scanner.Main()
	// return
	ASTPrinter := ast.NewASTPringer()
	fmt.Println(ASTPrinter.Print(testAST()))

	if len(args) > 2 {
		fmt.Println("usage: glox [script]")
	} else if len(args) == 2 {
		if args[1] == "test" {
			runFile(testInputFile)
		} else {
			runFile(args[1])
		}
		if hadError {
			os.Exit(65)
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
	hadError = false
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			run(line)
		}
		fmt.Print("> ")
	}
}

func run(i string) {
	s := scanner.NewScanner(i, *errReporter)
	ts := s.ScanTokens()
	for _, t := range ts {
		fmt.Printf("%s\n", t)
	}
}
