package main

import (
	"bufio"
	"fmt"
	"os"
)

func READ(input string) string {
	return input
}

func EVAL(input string) string {
	return input
}

func PRINT(input string) string {
	return input
}

func rep(input string) string {
	output := READ(input)
	output = EVAL(output)
	return PRINT(output)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("user> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			return
		}

		output := rep(input)
		fmt.Print(output)
	}
}
