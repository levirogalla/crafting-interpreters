//go:generate ./cmd/generate-ast/generate-ast ./ast/expr.grammar -o ./ast/expr.go
//go:generate ./cmd/generate-ast/generate-ast ./ast/stmt.grammar -o ./ast/stmt.go
package main

import (
	"bufio"
	"crafting-interpreters/ast"
	"crafting-interpreters/error"
	"crafting-interpreters/interpreter"
	"crafting-interpreters/models"
	"crafting-interpreters/parser"
	"crafting-interpreters/scanner"
	"fmt"
	"os"
)

const testInputFile = "/Users/levirogalla/Projects/lib/crafting-interpreters/main.lox"

var errReporter = error.NewReporter(false)

func testAST() ast.Expr {
	expr := &ast.BinaryNode{
		Left: ast.UnaryNode{
			Op: models.NewToken(models.Minus, "-", nil, 1),
			Right: ast.LiteralNode{
				Value: models.NewToken(models.Num, "123", 123, 1),
			},
		},
		Op: models.NewToken(models.Star, "*", nil, 1),
		Right: ast.GroupingNode{
			Expr: ast.LiteralNode{
				Value: models.NewToken(models.Num, "45.67", 45.67, 1),
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
func run(i string) {
	scanner := scanner.NewScanner(i, *errReporter)
	ts := scanner.ScanTokens()
	parser := parser.NewParser(ts)
	stmts, err := parser.Parse()
	_ = err
	interp.Interpret(stmts)
}
