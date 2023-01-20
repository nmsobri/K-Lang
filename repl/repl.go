package repl

import (
	"Klang/lexer"
	"Klang/token"
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

		tok := l.NextToken()
		for tok.Type != token.EOF {
			fmt.Println(tok)
			tok = l.NextToken()
		}
	}
}
