package handlers

import (
	"fmt"

	"github.com/abdelmounim-dev/redis/pkg/parser"
)

func HandleCommand(t *parser.Token) (*parser.Token, error) {
	if t.Type == parser.Array {
		switch string(t.Value.([]*parser.Token)[0].Value.([]byte)) {
		case "PING":
			return handlePing(t)

		}
	}
	return &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"}, nil
}

func handlePing(t *parser.Token) (*parser.Token, error) {
	if len(t.Value.([]*parser.Token)) == 1 {
		return &parser.Token{Type: parser.SimpleString, Value: "PONG"}, nil
	}
	fmt.Println(t)
	return t.Value.([]*parser.Token)[1], nil
}
