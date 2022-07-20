package main

import "testing"

func TestPeek(t *testing.T) {
	s := NewScanner("testing peek")
	got := s.peek()
	if got != 't' {
		t.Errorf("Peek return was incorrect, got %c, expected %c", got, 't')
	}

	s = NewScanner("")
	got = s.peek()
	if got != '\x00' {
		t.Errorf("Peek return was incorrect, got %c, expected %c", got, '\x00')
	}

	// Peek does not advance current
	if s.current != 0 {
		t.Errorf("Peek advanced 'current'")
	}

}

func TestIsAtEnd(t *testing.T) {
	s := NewScanner("test")
	got := s.isAtEnd()
	if got {
		t.Errorf("IsAtEnd returned true when not at end")
	}
	s.current = 50
	got = s.isAtEnd()
	if !got {
		t.Errorf("IsAtEnd returned false when at the end")
	}
}

func TestAdvance(t *testing.T) {
	s := NewScanner("test")
	got := s.advance()
	if got != 't' {
		t.Errorf("Advance return was incorrect, got %c, expected %c", got, 't')
	}
	if s.current != 1 {
		t.Errorf("Advance did not increase 'current' after returning")
	}
}

func TestAddTokenTypeObject(t *testing.T) {
	s := NewScanner("64")
	s.current = 2
	s.addTokenTypeObject(NUMBER, 64.0)
	expected := Token{NUMBER, "64", 64.0, 1}
	if s.tokens[0] != expected {
		t.Errorf("AddTokenTypeObject was incorrect, got {%+v}, expected {%+v}", s.tokens[0], expected)
	}
}

func TestAddToken(t *testing.T) {
	s := NewScanner("+")
	s.current = 1
	s.addToken(PLUS)
	expected := Token{PLUS, "+", nil, 1}
	if s.tokens[0] != expected {
		t.Errorf("AddToken was incorrect, got {%+v}, expected {%+v}", s.tokens[0], expected)
	}

	s = NewScanner("64")
	s.current = 2
	s.line = 5
	s.addToken(NUMBER, 64.0)
	expected = Token{NUMBER, "64", 64.0, 5}
	if s.tokens[0] != expected {
		t.Errorf("AddTokenTypeObject was incorrect, got {%+v}, expected {%+v}", s.tokens[0], expected)
	}
}