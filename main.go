package main

import (
	"charm/evaluator"
	"charm/lexer"
	"charm/object"
	"charm/parser"
	"charm/repl"
	"fmt"
	"os"
)

func main() {
    args := os.Args

    if len(args) == 1 {
        fmt.Println(args)
        fmt.Printf("Charm v0.1\n")
        repl.Start(os.Stdin, os.Stdout)
    } else if len(args) == 2 {
        filePath := args[1]
        file, err := os.ReadFile(filePath)
        if err != nil {
            fmt.Printf("error reading file: %s\n", filePath)
        }

        lexer := lexer.New(string(file))
        parser := parser.New(lexer)
        program := parser.ParseProgram()
        env := object.NewEnvironment()

        evaluator.Eval(program, env)
    } else {
        fmt.Println("incorrect number of arguments")
    }
}
