package main

import (
	"fmt"
	"strings"
	"unicode"
)

func lex(content string) ([]string, []string) {
	valid_tokens := []string{}
	invalid_tokens := []string{}
	current_line := 1	

	for i := 0; i < len(content); i++ {
		ch := content[i]

		switch ch {
			case '(':
				valid_tokens = append(valid_tokens, "LEFT_PAREN ( null")
			case ')':
				valid_tokens = append(valid_tokens, "RIGHT_PAREN ) null")
			case '{':
				valid_tokens = append(valid_tokens, "LEFT_BRACE { null")
			case '}':
				valid_tokens = append(valid_tokens, "RIGHT_BRACE } null")
			case ',':
				valid_tokens = append(valid_tokens, "COMMA , null")
			case '.':
				valid_tokens = append(valid_tokens, "DOT . null")
			case '-':
				valid_tokens = append(valid_tokens, "MINUS - null")
			case '+':
				valid_tokens = append(valid_tokens, "PLUS + null")
			case ';':
				valid_tokens = append(valid_tokens, "SEMICOLON ; null")
			case '*':
				valid_tokens = append(valid_tokens, "STAR * null")
			case '/':
				if i + 1 < len(content) && content[i + 1] == '/' {
					for i < len(content) && content[i] != '\n' {
						i++
					}
					i-- 
				} else {
					valid_tokens = append(valid_tokens, "SLASH / null")
				}
			case '=':
				if i + 1 < len(content) && content[i + 1] == '=' {
					valid_tokens = append(valid_tokens, "EQUAL_EQUAL == null")
					i += 1
				} else {
					valid_tokens = append(valid_tokens, "EQUAL = null")
				}
			case '!':
				if i + 1 < len(content) && content[i + 1] == '=' {
					valid_tokens = append(valid_tokens, "BANG_EQUAL != null")
					i += 1
				} else {
					valid_tokens = append(valid_tokens, "BANG ! null")
				}
			case '<':
				if i + 1 < len(content) && content[i + 1] == '=' {
					valid_tokens = append(valid_tokens, "LESS_EQUAL <= null")
					i += 1
				} else {
					valid_tokens = append(valid_tokens, "LESS < null")
				}
			case '>':
				if i + 1 < len(content) && content[i + 1] == '=' {
					valid_tokens = append(valid_tokens, "GREATER_EQUAL >= null")
					i += 1
				} else {
					valid_tokens = append(valid_tokens, "GREATER > null")
				}
			case '\n':
				current_line += 1
			case '\t', ' ':
				continue
			case '"':
				temp := ""
				valid_string := false
				i += 1
				
				for i < len(content) {
					if content[i] == '"' {
						valid_tokens = append(valid_tokens, fmt.Sprintf("STRING \"%s\" %s", temp, temp))
						valid_string = true
						break
					}
					if content[i] == '\n' {
						current_line += 1
					}
					temp += string(content[i])
					i += 1
				}
				
				if !valid_string {
					invalid_tokens = append(invalid_tokens, fmt.Sprintf("[line %d] Error: Unterminated string.", current_line))
				}
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				temp := ""
				dot := false
				
				for i < len(content) && (unicode.IsDigit(rune(content[i])) || content[i] == '.'){
					if content[i] == '.' {
						dot = true
					}
					temp += string(content[i])
					i += 1
				}
				i -= 1

				temp_parsed := temp
				if dot {
					parts := strings.Split(temp_parsed, ".")
					if len(parts) == 2 {
						allZeros := true
						for _, c := range parts[1] {
							if c != '0' {
								allZeros = false
								break
							}
						}
						if allZeros {
							temp_parsed = parts[0] + ".0"
						}
					}
				} else {
					temp_parsed = temp_parsed + ".0"
				}

				valid_tokens = append(valid_tokens, fmt.Sprintf("NUMBER %s %s", temp, temp_parsed))
			case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '_':
				temp := ""
				for i < len(content) && (unicode.IsLetter(rune(content[i])) || unicode.IsDigit(rune(content[i])) || content[i] == '_') {
					temp += string(content[i])
					i += 1
				}
				i -= 1

				if temp == "and" || temp == "class" || temp == "else" || temp == "false" || temp == "for" || temp == "fun" || temp == "if" || temp == "nil" || temp == "or" || temp == "print" || temp == "return" || temp == "super" || temp == "this" || temp == "true" || temp == "var" || temp == "while" {
					valid_tokens = append(valid_tokens, fmt.Sprintf("%s %s null", strings.ToUpper(temp), temp))
				} else {
					valid_tokens = append(valid_tokens, fmt.Sprintf("IDENTIFIER %s null", temp))
				}

			default:
				invalid_tokens = append(invalid_tokens, fmt.Sprintf("[line %d] Error: Unexpected character: %c", current_line, ch))

		}
	}

	return invalid_tokens, valid_tokens
}
