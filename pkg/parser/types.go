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

func (t *Token) Serialize() ([]byte, error) {
	fmt.Println(t)
	switch t.Type {
	case SimpleString:
		return []byte(fmt.Sprintf("+%s\r\n", t.Value.(string))), nil
	case SimpleError:
		return []byte(fmt.Sprintf("-%s\r\n", t.Value.(string))), nil
	case Integer:
		return []byte(fmt.Sprintf(":%d\r\n", t.Value.(int64))), nil
	case BulkString:
		if t.Value == nil {
			return []byte("$-1\r\n"), nil
		}
		bs := t.Value.([]byte)
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(bs), bs)), nil
	case Array:
		if t.Value == nil {
			return []byte("*-1\r\n"), nil
		}
		arr := t.Value.([]*Token)
		result := []byte(fmt.Sprintf("*%d\r\n", len(arr)))
		for _, elem := range arr {
			elemBytes, err := elem.Serialize()
			if err != nil {
				return nil, err
			}
			result = append(result, elemBytes...)
		}
		return result, nil
	case Null:
		return []byte("_\r\n"), nil
	case Boolean:
		if t.Value.(bool) {
			return []byte("#t\r\n"), nil
		}
		return []byte("#f\r\n"), nil
	case Double:
		return []byte(fmt.Sprintf(",%f\r\n", t.Value.(float64))), nil
	case BigNumber:
		return []byte(fmt.Sprintf("(%s\r\n", t.Value.(string))), nil
	case BulkError:
		bs := t.Value.([]byte)
		return []byte(fmt.Sprintf("!%d\r\n%s\r\n", len(bs), bs)), nil
	case VerbatimString:
		vs := t.Value.(struct{ Format, Text string })
		return []byte(fmt.Sprintf("=%d\r\n%s:%s\r\n", len(vs.Format)+1+len(vs.Text), vs.Format, vs.Text)), nil
	default:
		return nil, fmt.Errorf("type not implemented yet: %v", t.Type)
	}
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
