package repl

import (
	"Klang/lexer"
	"Klang/parser"
	"bufio"
	"fmt"
	"os"
)

const PROMPT = ">> "

func Start() {
	fmt.Println("Welcome To K Programming Language")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		input := scanner.Text()
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		for _, stmt := range program.Statements {
			fmt.Println(stmt.String())
		}
	}
}
