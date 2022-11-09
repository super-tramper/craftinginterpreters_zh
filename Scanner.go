package lox

import (
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token

	start   int
	current int
	line    int
}

func (s *Scanner) scanTokens() error {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() error {
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
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		err := s.toString()
		if err != nil {
			return err
		}
	case 'o':
		if s.peek() == 'r' {
			s.addToken(OR)
		}
	default:
		if s.isDigit(c) {
			s.toNumber()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			return &SyntaxError{s, "Unexpected character."}
		}
	}
	return nil
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

//func (s Scanner) addToken(typ TokenType, literal interface{}) {
//	s.addToken(typ, nil)
//	text := s.source[s.start, s.current]
//}

func (s *Scanner) addToken(typ TokenType) {
	s.addLiteralToken(typ, nil)
}

func (s *Scanner) addLiteralToken(typ TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	token := Token{typ, text, literal, s.line}
	s.tokens = append(s.tokens, token)
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

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\\'
	}
	return s.source[s.current]
}

func (s *Scanner) toString() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return &SyntaxError{s, "Unexpected character."}
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(STRING, value)
	return nil
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) toNumber() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	f, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		panic(err)
	}
	s.addLiteralToken(NUMBER, f)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\\'
	}
	return s.source[s.current+1]
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	typ, ok := keywords[text]
	if !ok {
		typ = IDENTIFIER
	}
	s.addToken(typ)
}
