package event

import "github.com/vmihailenco/msgpack/v4"

type LikeEvent struct {
	*Base
	UserID uint64
	Extra  string
}

func (o *LikeEvent) MarshalBinary() (data []byte, err error) {
	return msgpack.Marshal(o)
}

func (o *LikeEvent) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, o)
}
