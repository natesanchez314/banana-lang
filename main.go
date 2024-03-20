package main

import (
	"fmt"
	"os"
	"os/user"
	"banana/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to the Banana programming language!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}