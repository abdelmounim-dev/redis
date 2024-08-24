package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Token struct {
	tType DataType
	value interface{}
}

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
		return nil, err
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
	bytes, err := ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ SIMPLE STRING: %v", err)
	}
	t := &Token{
		tType: SimpleString,
		value: string(bytes),
	}
	return t, nil
}

func (p *Parser) readSimpleError() (*Token, error) {
	bytes, err := ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ SIMPLE ERROR: %v", err)
	}
	t := &Token{
		tType: SimpleError,
		value: string(bytes),
	}
	return t, nil
}

func (p *Parser) readInteger() (*Token, error) {
	bytes, err := ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ INTEGER: %v", err)
	}
	num, err := strconv.Atoi(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("READ INTEGER: %v", err)
	}
	t := &Token{
		tType: Integer,
		value: int64(num),
	}
	return t, nil
}

func (p *Parser) readBulkString() (*Token, error) {
	lenToken, err := p.readInteger()
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	n := lenToken.value.(int64)
	if n == 0 {
		_, err = ReadBytesUntilCRLF(p.buf)
		if err != nil {
			return nil, fmt.Errorf("READ BULK STRING: %v", err)
		}
		return &Token{tType: BulkString, value: []byte("")}, nil
	}
	if n == -1 {
		return &Token{tType: BulkString, value: nil}, nil
	}
	if n < -1 {
		return nil, fmt.Errorf("READ BULK STRING: Negative length (length < -1)")
	}
	b := make([]byte, n)
	_, err = io.ReadFull(p.buf, b)
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	_, err = ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ BULK STRING: %v", err)
	}
	t := &Token{
		tType: BulkString,
		value: b,
	}
	return t, nil
}

// TODO: implement
func (p *Parser) readArray() (*Token, error) {
	return nil, nil
}

func (p *Parser) readNull() (*Token, error) {
	_, err := ReadBytesUntilCRLF(p.buf)
	if err != nil {
		return nil, fmt.Errorf("READ NULL: %v", err)
	}
	return &Token{Null, nil}, nil
}

func (p *Parser) readBoolean() (*Token, error) {
	b, err := ReadBytesUntilCRLF(p.buf)
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
		tType: Boolean,
		value: v,
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
