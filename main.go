package main

import (
    "fmt"
    "os"
    "charm/repl"
)

func main() {
    fmt.Printf("Charm v0.1\n")
    repl.Start(os.Stdin, os.Stdout)
}
