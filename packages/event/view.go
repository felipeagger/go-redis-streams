package event

import "github.com/vmihailenco/msgpack/v4"

type ViewEvent struct {
	*Base
	UserID uint64
	Extra  string
}

func (o *ViewEvent) MarshalBinary() (data []byte, err error) {
	return msgpack.Marshal(o)
}

func (o *ViewEvent) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, o)
}
