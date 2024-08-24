package parser

import (
	"bufio"
	"io"
)

func ReadBytesUntilCRLF(reader *bufio.Reader) ([]byte, error) {
	var result []byte

	for {
		// Read until the next '\n'
		part, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF && len(part) > 0 {
				// Handle case where EOF is reached but data is still available
				result = append(result, part...)
				break
			}
			return nil, err
		}

		result = append(result, part...)

		// Check if the last two bytes are \r\n
		if len(result) >= 2 && result[len(result)-2] == '\r' && result[len(result)-1] == '\n' {
			return result[:len(result)-2], nil // Strip off the \r\n
		}
	}
	return result, nil
}
