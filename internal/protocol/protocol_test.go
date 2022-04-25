package protocol

import (
	"bytes"
	"testing"
)

func TestProtocol(t *testing.T) {
	var w = &bytes.Buffer{}
	t.Logf("err=%v", Encode(w, &Msg{Head: Head{Type: 0x11}, Content: []byte("wow")}))
	t.Logf("err=%+v", w)
	data, _ := Decode(w)
	t.Logf("content=%+v, type=%v", string(data.Content), data.Type)
}
