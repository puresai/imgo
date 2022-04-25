package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	VERSION = 001
	MAX_LEN = 10240
)

func Decode(r io.Reader) (*Msg, error) {
	var h Head

	// message size
	err := binary.Read(r, binary.BigEndian, &h)
	if err != nil {
		return nil, err
	}

	if h.Len > MAX_LEN {
		return nil, fmt.Errorf("response body size (%d) is greater than limit (%d)",
			h.Len, MAX_LEN)
	}

	m := &Msg{Head: h}

	// message binary data
	m.Content = make([]byte, h.Len)
	_, err = io.ReadFull(r, m.Content)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func Encode(w io.Writer, msg *Msg) error {
	var headBuff = make([]byte, 4)
	var big = binary.BigEndian
	headBuff[0] = msg.Head.Type
	big.PutUint16(headBuff[2:4], uint16(len(msg.Content)))

	if _, err := w.Write(headBuff); err != nil {
		return err
	}

	_, err := w.Write(msg.Content)
	if err != nil {
		return err
	}

	return nil
}
