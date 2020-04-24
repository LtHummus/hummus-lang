package main

import (
	"fmt"
	"hummus-lang/evaluator"
	"hummus-lang/lexer"
	"hummus-lang/object"
	"hummus-lang/parser"
	"hummus-lang/repl"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}

		l := lexer.New(string(contents))
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, e := range p.Errors() {
				fmt.Println(e)
				panic("invalid code")
			}
		}

		res := evaluator.Eval(program, object.NewEnvironment())
		if res.Type() == object.ErrorObj {
			fmt.Printf("%s\n", res.Printable())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Hello and welcome to the Hummus REPL.")
		fmt.Printf("\n")
		repl.Start(os.Stdin, os.Stdout)
	}

}
