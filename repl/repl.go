package repl

import (
	// "Klang/ast"
	"Klang/lexer"
	"Klang/parser"
	"bufio"
	"fmt"
	"log"
	"os"
)

const PROMPT = ">> "

func Start() {
	fmt.Println("Welcome To K Programming Language")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(PROMPT)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		input := scanner.Text()
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		for _, stmt := range program.Statements {
			// integ, _ := stmt.(*ast.ExpressionStatement)
			// fmt.Println(integ.Expression.String())
			fmt.Println(stmt.String())
		}
	}
}
