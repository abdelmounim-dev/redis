package parser

// DataType represents the different types of RESP data.
type DataType int

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
