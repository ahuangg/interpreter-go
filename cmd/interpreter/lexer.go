package main

import (
	"fmt"
)

type Token struct {
	Type    string
	Lexeme  string
	Literal interface{}
	Line    int
}

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	COMMA         = "COMMA"
	DOT           = "DOT"
	MINUS         = "MINUS"
	PLUS          = "PLUS"
	SEMICOLON     = "SEMICOLON"
	STAR          = "STAR"
	SLASH         = "SLASH"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG          = "BANG"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	STRING        = "STRING"
	NUMBER        = "NUMBER"
	IDENTIFIER    = "IDENTIFIER"
	EOF           = "EOF"
)

var keywords = map[string]string{
	"and":    "AND",
	"class":  "CLASS",
	"else":   "ELSE",
	"false":  "FALSE",
	"for":    "FOR",
	"fun":    "FUN",
	"if":     "IF",
	"nil":    "NIL",
	"or":     "OR",
	"print":  "PRINT",
	"return": "RETURN",
	"super":  "SUPER",
	"this":   "THIS",
	"true":   "TRUE",
	"var":    "VAR",
	"while":  "WHILE",
}

type Lexer struct {
	source  string
	tokens  []Token
	errors  []string
	start   int
	current int
	line    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		tokens: make([]Token, 0),
		errors: make([]string, 0),
		line:   1,
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() byte {
	l.current++
	return l.source[l.current-1]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) addToken(tokenType string, literal interface{}) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token{tokenType, text, literal, l.line})
}

func (l *Lexer) Tokenize() ([]Token, []string) {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, Token{EOF, "", nil, l.line})
	return l.tokens, l.errors
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case '(':
		l.addToken(LEFT_PAREN, nil)
	case ')':
		l.addToken(RIGHT_PAREN, nil)
	case '{':
		l.addToken(LEFT_BRACE, nil)
	case '}':
		l.addToken(RIGHT_BRACE, nil)
	case ',':
		l.addToken(COMMA, nil)
	case '.':
		l.addToken(DOT, nil)
	case '-':
		l.addToken(MINUS, nil)
	case '+':
		l.addToken(PLUS, nil)
	case ';':
		l.addToken(SEMICOLON, nil)
	case '*':
		l.addToken(STAR, nil)
	case '/':
		if l.match('/') {
			for l.peek() != '\n' && !l.isAtEnd() {
				l.advance()
			}
		} else {
			l.addToken(SLASH, nil)
		}
	case '=':
		if l.match('=') {
			l.addToken(EQUAL_EQUAL, nil)
		} else {
			l.addToken(EQUAL, nil)
		}
	case '!':
		if l.match('=') {
			l.addToken(BANG_EQUAL, nil)
		} else {
			l.addToken(BANG, nil)
		}
	case '<':
		if l.match('=') {
			l.addToken(LESS_EQUAL, nil)
		} else {
			l.addToken(LESS, nil)
		}
	case '>':
		if l.match('=') {
			l.addToken(GREATER_EQUAL, nil)
		} else {
			l.addToken(GREATER, nil)
		}
	case '\n':
		l.line++
	case ' ', '\r', '\t':
		break
	case '"':
		l.string()
	default:
		if isDigit(c) {
			l.number()
		} else if isAlpha(c) {
			l.identifier()
		} else {
			l.errors = append(l.errors, fmt.Sprintf("[line %d] Error: Unexpected character: %c", l.line, c))
		}
	}
}

func (l *Lexer) string() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		l.errors = append(l.errors, fmt.Sprintf("[line %d] Error: Unterminated string.", l.line))
		return
	}

	l.advance()
	value := l.source[l.start+1 : l.current-1]
	l.addToken(STRING, value)
}

func (l *Lexer) number() {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
	}

	value := l.source[l.start:l.current]
	l.addToken(NUMBER, value)
}

func (l *Lexer) identifier() {
	for isAlphaNumeric(l.peek()) {
		l.advance()
	}

	text := l.source[l.start:l.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	l.addToken(tokenType, nil)
}

func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}