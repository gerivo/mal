package main

import (
	"bufio"
	"fmt"
	"os"
)

func READ(input string) (MalTyper, error) {
	return read_str(input, true)
}

func EVAL(input MalTyper) MalTyper {
	return input
}

func PRINT(input MalTyper) string {
	return pr_str(input)
}

func rep(input string) string {
	output, error := READ(input)
	if error != nil {
		return error.Error()
	}

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
		fmt.Println(output)
	}
}
