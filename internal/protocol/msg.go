package protocol

type Head struct {
	Type uint8
	Todo uint8 // todo
	Len  uint16
}

type Msg struct {
	Head
	Content []byte
}
