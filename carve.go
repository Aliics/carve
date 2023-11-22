package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	program, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(program), "\n")

	instructs, err := parse(lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err = exec(defaultRuntime(), instructs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
