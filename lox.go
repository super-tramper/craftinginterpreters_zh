package lox

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
}

type SyntaxError struct {
	s       *Scanner
	message string
}

func (s *SyntaxError) Error() string {
	return report(s.s.line, "", s.message)
}

type EndOfContextError struct {
	s *Scanner
}

func (s *EndOfContextError) Error() string {
	return report(s.s.line, string(rune(s.s.current)), "end of context")
}

func report(line int, where string, message string) string {
	SetHadError(true)
	return fmt.Sprintf("[line%d] Error%s: %s", line, where, message)
}

func Run(args []string) {
	l := Lox{}
	if len(args) > 1 {
		fmt.Println("Usage: glox [script]")
		os.Exit(-1)
	} else if len(args) == 1 {
		l.runFile(args[0])
	} else {
		l.runPrompt()
	}
}

func (l *Lox) runFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	l.run(content)
	if GetHadError() {
		os.Exit(65)
	}
}

func (l *Lox) runPrompt() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		line, _, err := inputReader.ReadLine()
		if err != nil {
			fmt.Printf("ReadInput with err: %s\n", err)
		}
		if line == nil {
			break
		}
		l.run(line)
		SetHadError(false)
	}
}

//type Scanner struct {
//	source []byte
//}

//type Token string
//
//func (s Scanner) scanTokens() []Token {
//	var result []Token
//	for _, b := range s.source {
//		t := Token(b)
//		result = append(result, t)
//	}
//	return result
//}

func (l *Lox) run(source []byte) {
	scanner := Scanner{source: string(source)}
	err := scanner.scanTokens()
	if err != nil {
		panic(err)
	}
	for i, t := range scanner.tokens {
		fmt.Println(i, t)
	}
}
