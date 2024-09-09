package repl

import (
	"bufio"
	"charm/lexer"
    "charm/parser"
    "fmt"
    "io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()

        if !scanned {
            return
        }

        line := scanner.Text()
        lexer := lexer.New(line)
        parser := parser.New(lexer)
        program := parser.ParseProgram()

        errors := parser.GetErrors()
        if len(errors) != 0 {
            io.WriteString(out, "Parser errors:\n")
            for _, msg := range errors {
                io.WriteString(out, "\t" + msg + "\n")
            }
            continue
        }

        io.WriteString(out, program.String())
        io.WriteString(out, "\n")
    }
}
