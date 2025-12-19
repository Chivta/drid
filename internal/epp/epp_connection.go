package epp

import (
	"bufio"
	"net"
)

type Connection struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *Connection) ReadMessage() (string, error) {
	length, err := c.readMessageLength()
	if err != nil {
		return "", err
	}
	length -= 4

	message, err := c.readExact(length)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (c *Connection) WriteMessage(data []byte) error {
	totalLength := len(data) + 4
	lengthHeader := lengthToHeader(totalLength)

	_, err := c.writer.Write(lengthHeader)
	if err != nil {
		return err
	}
	_, err = c.writer.Write(data)
	if err != nil {
		return err
	}
	return c.writer.Flush()
}

func headerToLength(header []byte) int {
	return int(header[0])<<24 | int(header[1])<<16 | int(header[2])<<8 | int(header[3])
}

func lengthToHeader(length int) []byte {
	return []byte{
		byte((length >> 24) & 0xFF),
		byte((length >> 16) & 0xFF),
		byte((length >> 8) & 0xFF),
		byte(length & 0xFF),
	}
}

func (c *Connection) readMessageLength() (int, error) {
	header := make([]byte, 4)
	n, err := c.reader.Read(header)
	if err != nil || n != 4 {
		return 0, err
	}
	length := headerToLength(header)
	return length, nil
}

func (c *Connection) readExact(length int) ([]byte, error) {
	buf := make([]byte, length)
	totalRead := 0
	for totalRead < length {
		n, err := c.reader.Read(buf[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}
	return buf, nil
}
