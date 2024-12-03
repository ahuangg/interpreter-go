package main

import (
	"fmt"
	"os"
)

func print(invalid_tokens []string, valid_tokens []string) {
	for _, token := range invalid_tokens {
		fmt.Fprintln(os.Stderr, token)
	}

	for _, token := range valid_tokens {
		fmt.Println(token)
	}

	fmt.Println("EOF  null")
	
	if len(invalid_tokens) > 0 {
		os.Exit(65)
	}
}