package shared

import (
	"bufio"
	"io"
)

type SharedSocket interface {
	SendBytes([]byte)
	SendMessage(string)
}

func ReadBuffer(reader *bufio.Reader, size int) (int, []byte, error) {
	buffer := make([]byte, size)
	count, err := reader.Read(buffer)

	if count == 0 {
		if err == io.EOF {
			return 0, nil, nil
		} else if err != nil {
			return 0, nil, err
		}
	}

	return count, buffer, err
}
