package tokenizer

type Tokenizer struct {
	input        []byte
	position     int
	readPosition int
	ch           byte
}

func NewTokenizer(input []byte) *Tokenizer {
	t := &Tokenizer{input: input}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.readPosition >= len(t.input) {
		t.ch = 0
	} else {
		t.ch = t.input[t.readPosition]
	}
	t.position = t.readPosition
	t.readPosition += 1

}

func (t *Tokenizer) NextToken() Token {
	var tok Token
	var ttype TokenType
	switch t.ch {
	case '+':
		ttype = Plus
	case '-':
		ttype = Minus
	case ':':
		ttype = Colon
	case '$':
		ttype = Dollar
	case '*':
		ttype = Asterisk
	case '_':
		ttype = Underscore
	case '#':
		ttype = Hash
	case ',':
		ttype = Comma
	case '(':
		ttype = OpenParenthesis
	case '!':
		ttype = Exclamation
	case '=':
		ttype = Equals
	case '%':
		ttype = Percent
	case '~':
		ttype = Tilde
	case '>':
		ttype = GreaterThan
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		t.readChar()
		return tok
	default:

	}

	tok = NewToken(ttype, t.ch)
	t.readChar()

	return tok
}
