package tokenizer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// First Bytes
	Plus            = "+"
	Minus           = "-"
	Colon           = ":"
	Dollar          = "$"
	Asterisk        = "*"
	Underscore      = "_"
	Hash            = "#"
	Comma           = ","
	OpenParenthesis = "("
	Exclamation     = "!"
	Equals          = "="
	Percent         = "%"
	Tilde           = "~"
	GreaterThan     = ">"

	// Control Characters
	CRLF = "\r\n"

	// Aggregate type length
	LEN = "LENGTH"

	// General Error Handling
	Unknown = "unknown"
	Invalid = "invalid"

	// EOF
	EOF = ""
)

func NewToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
