package main

import (
	"fmt"
	"os"
	"os/user"
	"lemon/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to the Lemon programming language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}