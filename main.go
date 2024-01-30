package main

import(
    "fmt"
    "os"
    "A-Plus-Plus/repl"
)

func main() {
    fmt.Printf("Hello World!\n")
    fmt.Printf("Please enter commands\n")
    repl.Start(os.Stdin, os.Stdout)
}

