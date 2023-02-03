package main

import (
	"Klang/eval"
	"Klang/lexer"
	"Klang/object"
	"Klang/parser"
	"Klang/repl"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {

		contentBuff, err := ioutil.ReadFile(os.Args[1])
		content := string(contentBuff)

		if err != nil {
			log.Fatal(err)
		}

		env := object.NewEnvironment()
		l := lexer.New(content)
		p := parser.New(l)
		program := p.ParseProgram()

		eval.Eval(program, env)
	} else {
		repl.Start()
	}
}
