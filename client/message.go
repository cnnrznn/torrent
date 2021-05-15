package client

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type Message struct {
	Length  int
	ID      int
	Payload []byte
}

func NewMsg() ([]byte, error) {
	return nil, nil
}

func (c *Client) recvMsg(conn net.Conn) (*Message, error) {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, err
	}

	idBuf := make([]byte, 1)
	_, err = io.ReadFull(conn, idBuf)
	if err != nil {
		return nil, err
	}

	var len int
	var id int
	r := bytes.NewReader(lenBuf)
	err = binary.Read(r, binary.BigEndian, &len)
	if err != nil {
		return nil, err
	}

	r = bytes.NewReader(idBuf)
	err = binary.Read(r, binary.BigEndian, &id)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, len)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Length:  len,
		ID:      id,
		Payload: payload,
	}, nil
}

func (c *Client) handleMsg(msg Message) {
}

func (c *Client) sendInterested(conn net.Conn) {

}
