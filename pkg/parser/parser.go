package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/abdelmounim-dev/redis/pkg/utils"
)

type Parser struct {
	buf *bufio.Reader
}

func NewParser(buf *bufio.Reader) *Parser {
	p := &Parser{buf: buf}
	return p
}

func (p *Parser) NextToken() (*Token, error) {
	ch, err := p.buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("NEXT TOKEN: read Byte: %v", err)
	}
	switch ch {
	case '+':
		return p.readSimpleString()
	case '-':
		return p.readSimpleError()
	case ':':
		return p.readInteger()
	case '$':
		return p.readBulkString()
	case '*':
		return p.readArray()
	case '_':
		return p.readNull()
	case '#':
		return p.readBoolean()
	case ',':
		return p.readDouble()
	case '(':
		return p.readBigNumber()
	case '!':
		return p.readBulkError()
	case '=':
		return p.readVerbatimString()
	case '%':
		return p.readMap()
	case '~':
		return p.readSet()
	case '>':
		return p.readPush()
	default:
		return nil, fmt.Errorf("NEXT: Invalid First Byte (%v)", string(ch))
	}
}

func (p *Parser) readSimpleString() (*Token, error) {
	bytes, err := utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ SIMPLE STRING: %v", err)
	}
	t := &Token{
		Type:  SimpleString,
		Value: string(bytes),
	}
	return t, nil
}

func (p *Parser) readSimpleError() (*Token, error) {
	bytes, err := utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ SIMPLE ERROR: %v", err)
	}
	t := &Token{
		Type:  SimpleError,
		Value: string(bytes),
	}
	return t, nil
}

func (p *Parser) readInteger() (*Token, error) {
	bytes, err := utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ INTEGER: %v", err)
	}
	num, err := strconv.Atoi(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("READ INTEGER: %v", err)
	}
	t := &Token{
		Type:  Integer,
		Value: int64(num),
	}
	return t, nil
}

func (p *Parser) readBulkString() (*Token, error) {
	lenToken, err := p.readInteger()
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	n := lenToken.Value.(int64)
	if n == 0 {
		_, err = utils.ReadBytesUntilCRLF(p.buf)
		if err != nil {
			return nil, fmt.Errorf("READ BULK STRING: %v", err)
		}
		return &Token{Type: BulkString, Value: []byte("")}, nil
	}
	if n == -1 {
		return &Token{Type: BulkString, Value: nil}, nil
	}
	if n < -1 {
		return nil, fmt.Errorf("READ BULK STRING: Negative length (length < -1)")
	}
	b := make([]byte, n)
	_, err = io.ReadFull(p.buf, b)
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	_, err = utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	t := &Token{
		Type:  BulkString,
		Value: b,
	}
	return t, nil
}

// TODO: implement
func (p *Parser) readArray() (*Token, error) {
	lenToken, err := p.readInteger()
	if err != nil {
		return nil, fmt.Errorf("READ ARRAY: %v", err)
	}
	n := lenToken.Value.(int64)
	if n == 0 {
		_, err = utils.ReadBytesUntilCRLF(p.buf)
		if err != nil {
			return nil, fmt.Errorf("READ ARRAY: %v", err)
		}
		return &Token{Type: BulkString, Value: []byte("")}, nil
	}
	if n == -1 {
		return &Token{Type: BulkString, Value: nil}, nil
	}
	if n < -1 {
		return nil, fmt.Errorf("READ ARRAY: Negative length (length < -1)")
	}
	tokens := make([]*Token, 0, n)

	for i := int64(0); i < n; i++ {
		t, err := p.NextToken()
		if err != nil {
			return nil, fmt.Errorf("READ ARRAY: %v", err)
		}
		tokens = append(tokens, t)
	}

	t := &Token{Array, tokens}

	return t, nil
}

func (p *Parser) readNull() (*Token, error) {
	_, err := utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ NULL: %v", err)
	}
	return &Token{Null, nil}, nil
}

func (p *Parser) readBoolean() (*Token, error) {
	b, err := utils.ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ BOOLEAN: %v", err)
	}
	s := string(b)
	if s != "t" && s != "f" {
		return nil, fmt.Errorf("READ BOOLEAN: wrong token (neither 't' or 'f')")
	}
	var v bool
	if s == "t" {
		v = true
	} else {
		v = false
	}
	t := &Token{
		Type:  Boolean,
		Value: v,
	}
	return t, nil
}

func (p *Parser) readDouble() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readBigNumber() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readBulkError() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readVerbatimString() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readMap() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readSet() (*Token, error) {
	return nil, nil
}

// TODO: implement
func (p *Parser) readPush() (*Token, error) {
	return nil, nil
}
