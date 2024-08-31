package handlers

import (
	"reflect"
	"testing"

	"github.com/abdelmounim-dev/redis/pkg/parser"
)

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    *parser.Token
		expected *parser.Token
		wantErr  bool
	}{
		{
			name:     "Empty Command",
			input:    &parser.Token{Type: parser.Array, Value: []*parser.Token{}},
			expected: &parser.Token{Type: parser.SimpleError, Value: "Empty Command"},
			wantErr:  false,
		},
		{
			name: "PING Command",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("PING")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleString, Value: "PONG"},
			wantErr:  false,
		},
		{
			name: "PING Command with argument",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("PING")},
					{Type: parser.BulkString, Value: []byte("Hello")},
				},
			},
			expected: &parser.Token{Type: parser.BulkString, Value: []byte("Hello")},
			wantErr:  false,
		},
		{
			name: "ECHO Command",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("ECHO")},
					{Type: parser.BulkString, Value: []byte("Hello")},
				},
			},
			expected: &parser.Token{Type: parser.BulkString, Value: []byte("Hello")},
			wantErr:  false,
		},
		{
			name: "set Command",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("SET")},
					{Type: parser.BulkString, Value: []byte("name")},
					{Type: parser.BulkString, Value: []byte("AB")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleString, Value: "+OK\r\n"},
			wantErr:  false,
		},
		{
			name: "set Command with wrong type (integer)",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("SET")},
					{Type: parser.BulkString, Value: []byte("name")},
					{Type: parser.Integer, Value: int64(256)},
				},
			},
			expected: &parser.Token{Type: parser.SimpleError, Value: "-Wrong Type, IDK man\r\n"},
			wantErr:  false,
		},
		{
			name: "Unsupported Command",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("UNSUPPORTED")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"},
			wantErr:  false,
		},
		{
			name:     "Non-Array Input",
			input:    &parser.Token{Type: parser.SimpleString, Value: "Not an array"},
			expected: &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HandleCommand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HandleCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHandlePing(t *testing.T) {
	tests := []struct {
		name     string
		input    *parser.Token
		expected *parser.Token
		wantErr  bool
	}{
		{
			name: "PING without argument",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("PING")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleString, Value: "PONG"},
			wantErr:  false,
		},
		{
			name: "PING with argument",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("PING")},
					{Type: parser.BulkString, Value: []byte("Hello")},
				},
			},
			expected: &parser.Token{Type: parser.BulkString, Value: []byte("Hello")},
			wantErr:  false,
		},
		{
			name: "PING with too many arguments",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("PING")},
					{Type: parser.BulkString, Value: []byte("Hello")},
					{Type: parser.BulkString, Value: []byte("Extra")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleError, Value: "wrong argument number for PING command: 3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlePing(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("handlePing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHandleEcho(t *testing.T) {
	tests := []struct {
		name     string
		input    *parser.Token
		expected *parser.Token
		wantErr  bool
	}{
		{
			name: "ECHO with argument",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("ECHO")},
					{Type: parser.BulkString, Value: []byte("Hello")},
				},
			},
			expected: &parser.Token{Type: parser.BulkString, Value: []byte("Hello")},
			wantErr:  false,
		},
		{
			name: "ECHO with too many arguments",
			input: &parser.Token{
				Type: parser.Array,
				Value: []*parser.Token{
					{Type: parser.BulkString, Value: []byte("ECHO")},
					{Type: parser.BulkString, Value: []byte("Hello")},
					{Type: parser.BulkString, Value: []byte("Extra")},
				},
			},
			expected: &parser.Token{Type: parser.SimpleError, Value: "wrong argument number for PING command: 3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handleEcho(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleEcho() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("handleEcho() = %v, want %v", result, tt.expected)
			}
		})
	}
}
