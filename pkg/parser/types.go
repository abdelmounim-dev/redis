package parser

import "fmt"

// DataType represents the different types of RESP data.
type DataType int

// Token
type Token struct {
	Type  DataType
	Value interface{}
}

func (t Token) String() string {
	if t.Type == BulkString {
		return fmt.Sprintf("%v: %v", t.Type, string(t.Value.([]byte)))
	}
	return fmt.Sprintf("%v: %v", t.Type, t.Value)
}

// Enum values for DataType.
const (
	SimpleString DataType = iota
	SimpleError
	Integer
	BulkString
	Array
	Null
	Boolean
	Double
	BigNumber
	BulkError
	VerbatimString
	Map
	Set
	Push
)

// String returns the string representation of the DataType.
func (dt DataType) String() string {
	return [...]string{
		"SimpleString",
		"SimpleError",
		"Integer",
		"BulkString",
		"Array",
		"Null",
		"Boolean",
		"Double",
		"BigNumber",
		"BulkError",
		"VerbatimString",
		"Map",
		"Set",
		"Push",
	}[dt]
}

const (
	// First Bytes
	Plus            = "+"
	Minus           = "-"
	Colon           = ":"
	Dollar          = "$"
	Asterisk        = "*"
	Underscore      = "_"
	Hash            = "#"
	Comma           = ","
	OpenParenthesis = "("
	Exclamation     = "!"
	Equals          = "="
	Percent         = "%"
	Tilde           = "~"
	GreaterThan     = ">"

	// Control Characters
	CRLF = "\r\n"
)
