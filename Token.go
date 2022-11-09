package lox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) toString() string {
	return string(t.Type) + " " + t.Lexeme + " " + fmt.Sprintf("%v", t.Literal)
}
