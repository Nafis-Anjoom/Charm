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

        if !checkParserErrors(parser) {
            return
        }

        env := object.NewEnvironment()

        evaluator.Eval(program, env)
    } else {
        fmt.Println("incorrect number of arguments")
    }
}

func checkParserErrors(parser *parser.Parser) bool {
    errors := parser.GetErrors()
    if len(errors) == 0 {
        return true
    }

    fmt.Printf("parser has %d errors\n", len(errors))
    for _, msg := range errors {
        fmt.Printf("parser error: %q\n", msg)
    }

    return false
}
