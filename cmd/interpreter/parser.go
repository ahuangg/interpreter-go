package main

import (
	"fmt"
	"os"
)

type Expression interface {
	String() string
}

type Literal struct {
	Value interface{}
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

type Unary struct {
	Operator Token
	Right    Expression
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.Operator.Lexeme, u.Right.String())
}

type Binary struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Operator.Lexeme, b.Left.String(), b.Right.String())
}

type Grouping struct {
	Expression Expression
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.Expression.String())
}

type Parser struct {
	tokens      []Token
	cursor      int
	expressions []Expression
	errors      []error
}

type ParseError struct {
	Line    int
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:      tokens,
		cursor:      0,
		expressions: []Expression{},
		errors:      []error{},
	}
}

func (p *Parser) advance() Token {
	p.cursor++
	return p.tokens[p.cursor-1]
}

func (p *Parser) peek() Token {
	return p.tokens[p.cursor]
}

func (p *Parser) Parse() []Expression {
	for !p.isAtEnd() {
		expr := p.expression()
		if expr != nil {
			p.expressions = append(p.expressions, expr)
		}
	}
	return p.expressions
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == "EOF"
}

func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	expr := p.comparison()
	if expr == nil {
		return nil
	}

	for p.peek().Type == "EQUAL_EQUAL" || p.peek().Type == "BANG_EQUAL" {
		operator := p.advance()
		right := p.comparison()
		if right == nil {
			return nil
		}
		expr = Binary{Operator: operator, Left: expr, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()
	if expr == nil {
		return nil
	}

	for p.peek().Type == "GREATER" || p.peek().Type == "GREATER_EQUAL" ||
		p.peek().Type == "LESS" || p.peek().Type == "LESS_EQUAL" {
		operator := p.advance()
		right := p.term()
		if right == nil {
			return nil
		}
		expr = Binary{Operator: operator, Left: expr, Right: right}
	}

	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()
	if expr == nil {
		return nil
	}

	for p.peek().Type == "PLUS" || p.peek().Type == "MINUS" {
		operator := p.advance()
		right := p.factor()
		if right == nil {
			return nil
		}
		expr = Binary{Operator: operator, Left: expr, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()
	if expr == nil {
		return nil
	}

	for p.peek().Type == "STAR" || p.peek().Type == "SLASH" {
		operator := p.advance()
		right := p.unary()
		if right == nil {
			return nil
		}
		expr = Binary{Operator: operator, Left: expr, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expression {
	if p.peek().Type == "BANG" || p.peek().Type == "MINUS" {
		operator := p.advance()
		right := p.unary()
		if right == nil {
			return nil
		}
		return Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expression {
	token := p.advance()

	switch token.Type {
	case "TRUE", "FALSE", "NIL":
		return Literal{Value: token.Lexeme}
	case "STRING", "NUMBER":
		return Literal{Value: token.Literal}
	case "LEFT_PAREN":
		expr := p.expression()
		if expr == nil {
			return nil
		}
		if p.peek().Type != "RIGHT_PAREN" {
			return nil
		}
		p.advance()
		return Grouping{Expression: expr}
	default:
		p.errors = append(p.errors, ParseError{
			Line:    token.Line,
			Message: fmt.Sprintf("Error at '%s': Expect expression", token.Lexeme),
		})
		return nil
	}
}

func (p *Parser) Print() {
	if len(p.errors) > 0 {
		for _, err := range p.errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(65)
	}

	for _, exp := range p.expressions {
		fmt.Println(exp.String())
	}
}