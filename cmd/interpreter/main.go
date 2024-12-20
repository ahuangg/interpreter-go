package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if command == "tokenize" {
		lexer := NewLexer(string(fileContents))
		lexer.Tokenize()
		lexer.Print()
	}

	if command == "parse" {
		lexer := NewLexer(string(fileContents))
		valid_tokens, invalid_tokens := lexer.Tokenize()
		
		if len(invalid_tokens) > 0 {
			os.Exit(65)
		}

		parser := NewParser(valid_tokens)
		parser.Parse()
		parser.Print()


		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	os.Exit(1)
}
