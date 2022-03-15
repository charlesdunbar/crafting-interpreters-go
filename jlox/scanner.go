package main

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens(l *Lox) []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken(l)
	}
	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

func (s Scanner) scanToken(l *Lox) {
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
	default:
		l.error(s.line, "Unexpected character.")
	}
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

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s Scanner) addToken(l_type TokenType, literal ...interface{}) {
	if len(literal) == 0 {
		s.addTokenType(l_type)
	} else {
		s.addTokenTypeObject(l_type, literal[0])
	}
}

func (s Scanner) addTokenType(l_type TokenType) {
	s.addTokenTypeObject(l_type, nil)
}

func (s *Scanner) addTokenTypeObject(l_type TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{l_type, text, literal, s.line})
}
