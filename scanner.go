package main

import (
	"strconv"
)

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
		keywords: map[string]TokenType{
			"and":    AND,
			"class":  CLASS,
			"else":   ELSE,
			"false":  FALSE,
			"for":    FOR,
			"fun":    FUN,
			"if":     IF,
			"nil":    NIL,
			"or":     OR,
			"print":  PRINT,
			"return": RETURN,
			"super":  SUPER,
			"this":   THIS,
			"true":   TRUE,
			"var":    VAR,
			"while":  WHILE,
		},
	}
}

func (s *Scanner) ScanTokens(l *Lox) []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken(l)
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s *Scanner) scanToken(l *Lox) {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.multiComment(l)

		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++

	case '"':
		s.string(l)

	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			l.error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	if l_type, ok := s.keywords[text]; ok {
		s.addToken(l_type)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Scanner) isAlpha(c byte) bool {
	return ((c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c == '_'))
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

/*
	number adds a token of a number, including anything after a decimal place
*/
func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if (s.peek() == '.') && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	if n, err := strconv.ParseFloat(s.source[s.start:s.current], 64); err == nil {
		s.addToken(NUMBER, n)
	}

}

/*
	peekNext returns the byte of the next character,
	or a \0 if the next character is past the length of
	the line
*/
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) string(l *Lox) {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		l.error(s.line, "Unterminated string.")
		return
	}

	s.advance() // The closing ".

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
}

func (s *Scanner) multiComment(l *Lox) {
	for !(s.peek() == '*' && s.peekNext() == '/') && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		l.error(s.line, "Unterminated multiline comment.")
		return
	}

	s.advance() // The closing *.
	s.advance() // The closing /.
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

/*
	peek returns the next byte of the source,
	returning \0 if scanned past it

	it does not advance the Scanner's current position,
	like advance() does
*/
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

/*
	isAtEnd returns true if the Scanner's
	current value is greater than or equal
	to the length of the source
*/
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

/*
	advance returns the current byte, then advances
	the Scanner's current position.
*/
func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(l_type TokenType, literal ...any) {
	if len(literal) == 0 {
		s.addTokenTypeObject(l_type, nil)
	} else {
		s.addTokenTypeObject(l_type, literal[0])
	}
}

func (s *Scanner) addTokenTypeObject(l_type TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{l_type, text, literal, s.line})
}
