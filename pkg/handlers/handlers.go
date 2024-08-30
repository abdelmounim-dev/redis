package handlers

import (
	"strconv"

	"github.com/abdelmounim-dev/redis/pkg/parser"
)

func HandleCommand(t *parser.Token) (*parser.Token, error) {
	if t.Type == parser.Array {
		if len(t.Value.([]*parser.Token)) == 0 {
			return &parser.Token{Type: parser.SimpleError, Value: "Empty Command"}, nil
		}

		switch string(t.Value.([]*parser.Token)[0].Value.([]byte)) {
		case "PING":
			return handlePing(t)
		case "ECHO":
			return handleEcho(t)
		default:
			return &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"}, nil
		}
	}
	return &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"}, nil
}

func handlePing(t *parser.Token) (*parser.Token, error) {
	if len(t.Value.([]*parser.Token)) == 1 {
		return &parser.Token{Type: parser.SimpleString, Value: "PONG"}, nil
	}
	if len(t.Value.([]*parser.Token)) > 2 {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "wrong argument number for PING command: " + strconv.Itoa(len(t.Value.([]*parser.Token))),
		}, nil
	}
	return t.Value.([]*parser.Token)[1], nil
}
func handleEcho(t *parser.Token) (*parser.Token, error) {
	if len(t.Value.([]*parser.Token)) > 2 {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "wrong argument number for PING command: " + strconv.Itoa(len(t.Value.([]*parser.Token))),
		}, nil
	}
	return t.Value.([]*parser.Token)[1], nil
}
