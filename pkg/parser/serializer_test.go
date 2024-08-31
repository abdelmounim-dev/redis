package parser

import (
	"bytes"
	"testing"
)

func TestTokenSerialize(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected []byte
		wantErr  bool
	}{
		{
			name:     "SimpleString",
			token:    Token{Type: SimpleString, Value: "OK"},
			expected: []byte("+OK\r\n"),
		},
		{
			name:     "SimpleError",
			token:    Token{Type: SimpleError, Value: "Error message"},
			expected: []byte("-Error message\r\n"),
		},
		{
			name:     "Integer",
			token:    Token{Type: Integer, Value: int64(42)},
			expected: []byte(":42\r\n"),
		},
		{
			name:     "BulkString",
			token:    Token{Type: BulkString, Value: []byte("hello")},
			expected: []byte("$5\r\nhello\r\n"),
		},
		{
			name:     "NullBulkString",
			token:    Token{Type: BulkString, Value: nil},
			expected: []byte("$-1\r\n"),
		},
		{
			name: "Array",
			token: Token{Type: Array, Value: []*Token{
				{Type: SimpleString, Value: "a"},
				{Type: Integer, Value: int64(1)},
			}},
			expected: []byte("*2\r\n+a\r\n:1\r\n"),
		},
		{
			name:     "NullArray",
			token:    Token{Type: Array, Value: nil},
			expected: []byte("*-1\r\n"),
		},
		{
			name:     "Null",
			token:    Token{Type: Null},
			expected: []byte("_\r\n"),
		},
		{
			name:     "BooleanTrue",
			token:    Token{Type: Boolean, Value: true},
			expected: []byte("#t\r\n"),
		},
		{
			name:     "BooleanFalse",
			token:    Token{Type: Boolean, Value: false},
			expected: []byte("#f\r\n"),
		},
		{
			name:     "Double",
			token:    Token{Type: Double, Value: 3.14},
			expected: []byte(",3.140000\r\n"),
		},
		{
			name:     "BigNumber",
			token:    Token{Type: BigNumber, Value: "12345678901234567890"},
			expected: []byte("(12345678901234567890\r\n"),
		},
		{
			name:     "BulkError",
			token:    Token{Type: BulkError, Value: []byte("Error details")},
			expected: []byte("!13\r\nError details\r\n"),
		},
		{
			name: "VerbatimString",
			token: Token{Type: VerbatimString, Value: struct{ Format, Text string }{
				Format: "txt",
				Text:   "Some text",
			}},
			expected: []byte("=13\r\ntxt:Some text\r\n"),
		},
		{
			name:    "UnsupportedType",
			token:   Token{Type: 999},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.token.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Token.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(got, tt.expected) {
				t.Errorf("Token.Serialize() = %v, want %v", got, tt.expected)
			}
		})
	}
}
