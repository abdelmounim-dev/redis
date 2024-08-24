package parser

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestRESPParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Token
		wantErr  bool
	}{
		{"NullBulkString", "$-1\r\n", &Token{tType: BulkString, value: nil}, false},
		{"SimpleString", "+OK\r\n", &Token{tType: SimpleString, value: "OK"}, false},
		{"Error", "-Error message\r\n", &Token{tType: SimpleError, value: "Error message"}, false},
		{"Integer", ":1000\r\n", &Token{tType: Integer, value: int64(1000)}, false},
		{"BulkString", "$5\r\nhello\r\n", &Token{tType: BulkString, value: []byte("hello")}, false},
		// {"Array", "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", []interface{}{"hello", "world"}, false},
		// {"NestedArray", "*2\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$6\r\nfoobar\r\n", []interface{}{[]interface{}{"hello", "world"}, "foobar"}, false},
		// {"PingCommand", "*1\r\n$4\r\nping\r\n", []interface{}{"ping"}, false},
		// {"EchoCommand", "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n", []interface{}{"echo", "hello world"}, false},
		// {"GetCommand", "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", []interface{}{"get", "key"}, false},
		//
		// // Invalid cases
		// {"InvalidBulkStringLength", "$5\r\nhell\r\n", nil, true},
		// {"InvalidArrayLength", "*2\r\n$5\r\nhello\r\n", nil, true},
		// {"MissingCRLF", "$5\r\nhello", nil, true},
		// {"InvalidDataType", "x5\r\nhello\r\n", nil, true},
		// {"EmptyInput", "", nil, true},
		// {"IncompleteArray", "*2\r\n$5\r\nhello\r\n$5\r\nworl", nil, true},
	}

	b := bytes.NewBuffer(nil) // Create an empty buffer
	w := bufio.NewWriter(b)
	r := bufio.NewReader(b)
	p := NewParser(r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b.Reset()
			w.Reset(b)

			// Write the test input to the buffer
			_, err := w.WriteString(tt.input)
			if err != nil {
				t.Fatalf("Failed to write to buffer: %v", err)
			}

			// Flush the writer to ensure all data is written to the buffer
			err = w.Flush()
			if err != nil {
				t.Fatalf("Failed to flush writer: %v", err)
			}

			result, err := p.NextToken()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Parse() = %v, want %v", result, tt.expected)
			}
		})
	}
}
