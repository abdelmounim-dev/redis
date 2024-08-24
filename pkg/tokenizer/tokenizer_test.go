package tokenizer

import "testing"

func TestNextToken(t *testing.T) {
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{Plus, "+"},
		{Minus, "-"},
		{Colon, ":"},
		{Dollar, "$"},
		{Asterisk, "*"},
		{Underscore, "_"},
		{Hash, "#"},
		{Comma, ","},
		{OpenParenthesis, "("},
		{Exclamation, "!"},
		{Equals, "="},
		{Percent, "%"},
		{Tilde, "~"},
		{GreaterThan, ">"},
	}

	input := "+-:$*_#,(!=%~>"
	tokenizer := NewTokenizer([]byte(input))

	for i, tt := range tests {
		token := tokenizer.NextToken() // This function would return the next token based on the input.

		if token.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}

}
