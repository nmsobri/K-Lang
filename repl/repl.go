package repl

import (
	"Klang/environment"
	"Klang/eval"
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
	env := environment.New()

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

		evaluated := eval.Eval(program, env)
		fmt.Println(evaluated.Inspect())
	}
}
