package lexer

import "strings"

type Lexer struct {
	input string
	pos   int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Pos: l.pos}
	}

	start := l.pos

	if strings.HasPrefix(l.input[start:], "{{") {
		return l.readUntil("}}", TokenVariable)
	} else if strings.HasPrefix(l.input[start:], "{%") {
		return l.readUntil("%}", TokenBlock)
	} else if strings.HasPrefix(l.input[start:], "{#") {
		return l.readUntil("#}", TokenComment)
	} else {
		return l.readText()
	}
}

func (l *Lexer) readUntil(endTag string, kind TokenType) Token {
	end := strings.Index(l.input[l.pos:], endTag)
	if end == -1 {
		tok := Token{Type: kind, Literal: l.input[l.pos:], Pos: l.pos}
		l.pos = len(l.input)
		return tok
	}
	end += l.pos
	value := l.input[l.pos+2 : end]
	tok := Token{Type: kind, Literal: strings.TrimSpace(value), Pos: l.pos}
	l.pos = end + len(endTag)
	return tok
}

func (l *Lexer) readText() Token {
	start := l.pos
	for l.pos < len(l.input) &&
		!strings.HasPrefix(l.input[l.pos:], "{{") &&
		!strings.HasPrefix(l.input[l.pos:], "{%") &&
		!strings.HasPrefix(l.input[l.pos:], "{#") {
		l.pos++
	}
	return Token{Type: TokenText, Literal: l.input[start:l.pos], Pos: start}
}
