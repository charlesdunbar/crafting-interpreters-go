package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
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
	}
}

func (p *Parser) parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		dec, err := p.declaration()

		if err != nil {
			break
		}

		statements = append(statements, dec)
	}
	return statements
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}

	state, err := p.statement()

	if err != nil {
		p.synchronize()
		return nil, err
	}
	return state, nil
}

func (p *Parser) classDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect class name")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(LEFT_BRACE, "Expect '{' before class body.")
	if err != nil {
		return nil, err
	}

	var methods []Function
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		meth, err := p.function("method")

		if err != nil {
			return nil, err
		}
		methods = append(methods, *meth)
	}
	_, err = p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	if err != nil {
		return nil, err
	}
	return &Class{name: name, methods: methods}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		blo, err := p.block()
		if err != nil {
			return nil, err
		}
		return &Block{blo}, nil
	}
	return p.expressionStatement()

}

func (p *Parser) forStatement() (Stmt, error) {
	// Desugar a for-loop to a while loop
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// Initializer
	var initializer Stmt = nil
	if p.match(SEMICOLON) {
		// Initalizer skipped in for loop
		initializer = nil
	} else if p.match(VAR) {
		// Declare a variable scoped to for loop
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		// Init variable already exists
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	// Condition
	var condition Expr = nil
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	// Increment
	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &Block{[]Stmt{body, &Expression{increment}}}
	}

	if condition == nil {
		condition = &Literal{true}
	}
	body = &While{condition, body}

	if initializer != nil {
		body = &Block{[]Stmt{initializer, body}}
	}

	return body, nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &If{condition, thenBranch, elseBranch}, nil

}

func (p *Parser) printStatement() (Stmt, error) {
	value, _ := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after value.")

	if err != nil {
		return nil, err
	}

	return &Print{value}, nil
}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	var value Expr
	var err error
	// If the next thing after 'return' is a semicolon, it can't be an expression
	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &Return{keyword, value}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	var initial Expr
	name, err := p.consume(IDENTIFIER, "Expect variable name.")

	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		initial, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration")

	if err != nil {
		return nil, err
	}

	return &Var{name, initial}, nil

}

func (p *Parser) whileStatement() (Stmt, error) {
	_, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &While{condition, body}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, _ := p.expression()
	_, err := p.consume(SEMICOLON, "Expect ';' after expression.")

	if err != nil {
		return nil, err
	}

	return &Expression{expr}, nil
}

/*
 * Need to use *Function to return a nil if error occurs
 */
func (p *Parser) function(kind string) (*Function, error) {
	name, err := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if err != nil {
		return nil, err
	}
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	var params []Token
	if !p.check(RIGHT_PAREN) {
		// Do-while loop
		for {
			if len(params) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters.")
			}
			to_append, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, to_append)

			// Condition part of do-while
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return &Function{name, params, body}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements, nil
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if v, ok := expr.(*Variable); ok {
			name := v.name
			return &Assign{name, value}, nil
		} else if g, ok := expr.(*Get); ok {
			return &Set{g.object, g.name, value}, nil
		}
		return nil, p.error(equals, "Invalid assignment target.")
	}
	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &Logical{expr, operator, right}
	}

	return expr, nil

}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &Logical{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		expr = &Binary{expr, operator, right}
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
		return &Unary{operator, right}, nil
	}
	ret, err := p.call()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

/*
finishCall checks for any arguments passed to a function and calls itself for each argument sent
If there are no arguments we don't try to parse
*/
func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		// Mimic do-while loop
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			exp, err := p.expression()
			if err != nil {
				fmt.Println("Error in finish Call")
			}
			arguments = append(arguments, exp)

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		fmt.Println("Error in finish Call consume")
	}

	return &Call{callee, paren, arguments}

}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			expr = &Get{expr, name}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &Literal{false}, nil
	}
	if p.match(TRUE) {
		return &Literal{true}, nil
	}
	if p.match(NIL) {
		return &Literal{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}, nil
	}

	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{expr}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
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
	return Token{}, p.error(p.peek(), message)
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

func (p *Parser) error(token Token, message string) error {
	tokenError(token, message)
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
