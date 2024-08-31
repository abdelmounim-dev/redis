package handlers

import (
	"strconv"
	"strings"

	"github.com/abdelmounim-dev/redis/pkg/parser"
	"github.com/abdelmounim-dev/redis/pkg/storage"
)

type Handlers struct {
	store storage.Store
}

func NewHandlers(store storage.Store) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) HandleCommand(t *parser.Token) (*parser.Token, error) {
	if t.Type == parser.Array {
		if len(t.Value.([]*parser.Token)) == 0 {
			return &parser.Token{Type: parser.SimpleError, Value: "Empty Command"}, nil
		}

		key := strings.ToLower(string(t.Value.([]*parser.Token)[0].Value.([]byte)))

		switch key {
		case "ping":
			return h.handlePing(t)
		case "echo":
			return h.handleEcho(t)
		case "set":
			return h.handleSet(t)
		case "get":
			return h.handleGet(t)
		default:
			return &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"}, nil
		}
	}
	return &parser.Token{Type: parser.SimpleError, Value: "Command Not Supported"}, nil
}

func (h *Handlers) handlePing(t *parser.Token) (*parser.Token, error) {
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
func (h *Handlers) handleEcho(t *parser.Token) (*parser.Token, error) {
	if len(t.Value.([]*parser.Token)) != 2 {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "wrong argument number for ECHO command: " + strconv.Itoa(len(t.Value.([]*parser.Token))),
		}, nil
	}
	return t.Value.([]*parser.Token)[1], nil
}

func (h *Handlers) handleSet(t *parser.Token) (*parser.Token, error) {

	if len(t.Value.([]*parser.Token)) != 3 {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "wrong argument number for PING command: " + strconv.Itoa(len(t.Value.([]*parser.Token))),
		}, nil
	}

	if t.Value.([]*parser.Token)[1].Type != parser.BulkString {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "Wrong Key Type: " + t.Value.([]*parser.Token)[1].Type.String(),
		}, nil
	}

	if t.Value.([]*parser.Token)[2].Type != parser.BulkString {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "Wrong Value Type: " + t.Value.([]*parser.Token)[2].Type.String(),
		}, nil
	}

	key := string(t.Value.([]*parser.Token)[1].Value.([]byte))
	value := t.Value.([]*parser.Token)[2]

	err := h.store.Set(key, value)
	if err != nil {
		return nil, err
	}

	return &parser.Token{Type: parser.SimpleString, Value: "OK"}, nil
}
func (h *Handlers) handleGet(t *parser.Token) (*parser.Token, error) {
	if len(t.Value.([]*parser.Token)) != 2 {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "wrong argument number for PING command: " + strconv.Itoa(len(t.Value.([]*parser.Token))),
		}, nil
	}

	if t.Value.([]*parser.Token)[1].Type != parser.BulkString {
		return &parser.Token{
			Type:  parser.SimpleError,
			Value: "Wrong Key Type: " + t.Value.([]*parser.Token)[1].Type.String(),
		}, nil
	}
	key := string(t.Value.([]*parser.Token)[1].Value.([]byte))

	value, err := h.store.Get(key)
	if err != nil {
		return nil, err
	}

	return value, nil
}
