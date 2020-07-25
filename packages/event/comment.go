package event

import "github.com/vmihailenco/msgpack/v4"

type CommentEvent struct {
	*Base
	UserID  uint64
	Comment string
}

func (o *CommentEvent) MarshalBinary() (data []byte, err error) {
	return msgpack.Marshal(o)
}

func (o *CommentEvent) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, o)
}
