package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
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
		{"EmptyBulkString", "$0\r\n\r\n", &Token{tType: BulkString, value: []byte("")}, false},
		{"Boolean True", "#t\r\n", &Token{Boolean, true}, false},
		{"Boolean False", "#f\r\n", &Token{Boolean, false}, false},
		{"Array", "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n", &Token{Array, []*Token{{BulkString, []byte("hello")}, {BulkString, []byte("world")}}}, false},
		{"NestedArray", "*2\r\n*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n$6\r\nfoobar\r\n", &Token{Array, []*Token{{Array, []*Token{{BulkString, []byte("hello")}, {BulkString, []byte("world")}}}, {BulkString, []byte("foobar")}}}, false},
		{"PingCommand", "*1\r\n$4\r\nping\r\n", &Token{Array, []*Token{{BulkString, []byte("ping")}}}, false},
		{"EchoCommand", "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n", &Token{Array, []*Token{{BulkString, []byte("echo")}, {BulkString, []byte("hello world")}}}, false},
		{"GetCommand", "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", &Token{Array, []*Token{{BulkString, []byte("get")}, {BulkString, []byte("key")}}}, false},
		//
		// // Invalid cases
		{"InvalidBulkStringLength", "$5\r\nhell\r\n", nil, true},
		{"InvalidArrayLength", "*2\r\n$5\r\nhello\r\n", nil, true},
		{"MissingCRLF", "$5\r\nhello", nil, true},
		{"InvalidDataType", "x5\r\nhello\r\n", nil, true},
		{"EmptyInput", "", nil, true},
		{"IncompleteArray", "*2\r\n$5\r\nhello\r\n$5\r\nworl", nil, true},
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

			if !deepEqualTokens(result, tt.expected) {
				t.Errorf("Parse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// deepEqualTokens compares two Tokens for deep equality
func deepEqualTokens(a, b *Token) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.tType != b.tType {
		return false
	}
	switch a.tType {
	case Array:
		aArr, aOk := a.value.([]*Token)
		bArr, bOk := b.value.([]*Token)
		if !aOk || !bOk {
			fmt.Printf("Array type assertion failed: aOk = %v, bOk = %v\n", aOk, bOk)
			fmt.Printf("a.value type: %T, b.value type: %T\n", a.value, b.value)
			return false
		}
		if len(aArr) != len(bArr) {
			fmt.Println(2)
			return false
		}
		for i := range aArr {
			if !deepEqualTokens(aArr[i], bArr[i]) {
				fmt.Println(3)
				return false
			}
		}
		return true
	// case BulkString:
	// 	aStr, aOk := a.value.([]byte)
	// 	bStr, bOk := b.value.([]byte)
	// 	return aOk && bOk && bytes.Equal(aStr, bStr)
	// Add cases for other types as needed
	default:
		return reflect.DeepEqual(a.value, b.value)
	}
}

// tokenToString converts a Token to a string representation for error messages
func tokenToString(t *Token) string {
	if t == nil {
		return "<nil>"
	}
	switch t.tType {
	case Array:
		arr, ok := t.value.([]*Token)
		if !ok {
			return "<invalid array>"
		}
		elements := make([]string, len(arr))
		for i, elem := range arr {
			elements[i] = tokenToString(elem)
		}
		return fmt.Sprintf("Array[%s]", strings.Join(elements, ", "))
	case BulkString:
		str, ok := t.value.([]byte)
		if !ok {
			return "<invalid bulk string>"
		}
		return fmt.Sprintf("BulkString[%s]", string(str))
	// Add cases for other types as needed
	default:
		return fmt.Sprintf("%v", t.value)
	}
}
