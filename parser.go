package main

import "fmt"

type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
}

type ParseError struct {
	err error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("ParseError: %v", e.err)
}

func NewParser(tokens []Token, l *Lox) *Parser {
	return &Parser{
		current: 0,
		tokens:  tokens,
		lox:     l,
	}
}

func (p *Parser) parse() Expr {
	e, err := p.expression()
	if err != nil {
		return nil
	}
	return e
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()

	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()

	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()

	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{operator, right}, nil
	}
	x, err := p.primary()
	if err != nil {
		return nil, err
	}
	return x, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return Literal{false}, nil
	}
	if p.match(TRUE) {
		return Literal{true}, nil
	}
	if p.match(NIL) {
		return Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return Literal{p.previous().literal}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.", p.lox)
}

func (p *Parser) match(types ...TokenType) bool {
	for _, l_type := range types {
		if p.check(l_type) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(l_type TokenType, message string) (Token, error) {
	if p.check(l_type) {
		return p.advance(), nil
	}
	return Token{}, p.error(p.peek(), message, p.lox)
}

func (p *Parser) check(l_type TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().l_type == l_type
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().l_type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token Token, message string, l *Lox) error {
	l.tokenError(token, message)
	return ParseError{}
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().l_type == SEMICOLON {
			return
		}
		switch p.peek().l_type {
		case CLASS, FOR, FUN, IF, PRINT, RETURN, VAR, WHILE:
			return
		}
		p.advance()
	}
}
