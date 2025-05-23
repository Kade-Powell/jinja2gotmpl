package lexer

type TokenType string

const (
	TokenText     TokenType = "TEXT"
	TokenVariable TokenType = "VARIABLE"
	TokenBlock    TokenType = "BLOCK"
	TokenComment  TokenType = "COMMENT"
	TokenEOF      TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}
