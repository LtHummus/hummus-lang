package main

import (
	"fmt"
	"hummus-lang/repl"
	"os"
)

func main() {
	fmt.Printf("Hello and welcome to the Hummus REPL.")
	fmt.Printf("\n")
	repl.Start(os.Stdin, os.Stdout)
}
