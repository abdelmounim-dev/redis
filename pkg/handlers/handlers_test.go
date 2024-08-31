package handlers

import (
	"reflect"
	"testing"

	"github.com/abdelmounim-dev/redis/pkg/parser"
	"github.com/abdelmounim-dev/redis/pkg/storage"
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
			expected: &parser.Token{Type: parser.SimpleString, Value: "OK"},
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
			expected: &parser.Token{Type: parser.SimpleError, Value: "Wrong Value Type: " + parser.Integer.String()},
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
			s := storage.NewKeyValueStore()
			h := NewHandlers(s)
			result, err := h.HandleCommand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("name: %v, HandleCommand() = %v, want %v", tt.name, result, tt.expected)
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
			s := storage.NewKeyValueStore()
			h := NewHandlers(s)
			result, err := h.handlePing(tt.input)
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
				},
			},
			expected: &parser.Token{Type: parser.SimpleError, Value: "wrong argument number for ECHO command: 1"},
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
			expected: &parser.Token{Type: parser.SimpleError, Value: "wrong argument number for ECHO command: 3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewKeyValueStore()
			h := NewHandlers(s)
			result, err := h.handleEcho(tt.input)
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

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(s storage.Store)
		input     string
		expected  *parser.Token
		expectErr bool
	}{
		{
			name: "Existing Key",
			setup: func(s storage.Store) {
				s.Set(
					"name",
					&parser.Token{Type: parser.BulkString, Value: []byte("ab")},
				)
			},
			input:     "name",
			expected:  &parser.Token{Type: parser.BulkString, Value: []byte("ab")},
			expectErr: false,
		},
		{
			name: "Non-Existing Key",
			setup: func(s storage.Store) {
				// No setup, key doesn't exist
			},
			input:     "nonexistent",
			expected:  &parser.Token{Type: parser.BulkString, Value: nil},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewKeyValueStore() // Initialize your store
			tt.setup(s)                     // Run setup function to populate data

			result, err := s.Get(tt.input) // Call the Get method

			if (err != nil) != tt.expectErr {
				t.Errorf("Get() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("name: %v, Get() = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}
