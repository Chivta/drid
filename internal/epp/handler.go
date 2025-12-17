package epp

import (
	"bufio"
	"log"
	"net"
)

type Handler interface {
	Handle(conn net.Conn)
}

func NewEPPHandler() Handler {
	return &eppHandler{}
}

type eppHandler struct {
	reader *bufio.Reader
	writer *bufio.Writer
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

func (h *eppHandler) readMessage() (string, error){
	length, err := h.readMessageLength()
	if err != nil {
		return "", err
	}
	length -= 4 // subtract header size

	message, err := h.readExact(length)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

func (h *eppHandler) readMessageLength() (int, error){
	header := make([]byte, 4)
	n, err := h.reader.Read(header)
	if err != nil || n != 4 {
		return 0, err
	}
	length := headerToLength(header)
	return length, nil
}

func (h *eppHandler) readExact(length int) ([]byte, error) {
	buf := make([]byte, length)
	totalRead := 0
	for totalRead < length {
		n, err := h.reader.Read(buf[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}
	return buf, nil
}

func (h *eppHandler) writeMessage(data []byte) error {
	totalLength := len(data) + 4
	lengthHeader := lengthToHeader(totalLength)

	_, err := h.writer.Write(lengthHeader)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(data)
	if err != nil {
		return err
	}
	return h.writer.Flush()
}

const greetingResponse = `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <greeting>
    <svID>UANIC EPP-TEST2 Server version 0.6.3-4</svID>
    <svDate>2025-12-17T13:00:41.512+02:00</svDate>
    <svcMenu>
      <version>1.0</version>
      <lang>en</lang>
      <lang>ru</lang>
      <lang>ua</lang>
      <objURI>urn:ietf:params:xml:ns:epp-1.0</objURI>
      <objURI>http://hostmaster.ua/epp/contact-1.1</objURI>
      <objURI>http://hostmaster.ua/epp/domain-1.1</objURI>
      <objURI>http://hostmaster.ua/epp/host-1.1</objURI>
      <svcExtension>
        <extURI>http://hostmaster.ua/epp/rgp-1.1</extURI>
        <extURI>http://hostmaster.ua/epp/uaepp-1.1</extURI>
        <extURI>http://hostmaster.ua/epp/balance-1.0</extURI>
        <extURI>http://hostmaster.ua/epp/secDNS-1.1</extURI>
      </svcExtension>
    </svcMenu>
    <dcp>
      <access>
        <all/>
      </access>
      <statement>
        <purpose>
          <admin/>
          <prov/>
        </purpose>
        <recipient>
          <public/>
        </recipient>
        <retention>
          <stated/>
        </retention>
      </statement>
    </dcp>
  </greeting>
</epp>`

type Client struct {
	ClientId	  string
	Authenticated bool
}

func (h *eppHandler) Handle(conn net.Conn) {
	defer conn.Close()

	h.reader = bufio.NewReader(conn)
	h.writer = bufio.NewWriter(conn)

	client := &Client{
		ClientId:      "",
		Authenticated: false,
	}

	response := []byte(greetingResponse)

	err := h.writeMessage(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	log.Printf("Sent greeting response")

	for {
		request, err := h.readMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		response, err := HandleRequest(client, request)
		if err != nil {
			log.Printf("Error handling request: %v", err)
			return
		}
		err = h.writeMessage([]byte(response))
		if err != nil {
			log.Printf("Error writing response: %v", err)
			return
		}
		log.Printf("Sent response")
	}
}