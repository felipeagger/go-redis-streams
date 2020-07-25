package event

import (
	"encoding"
	"fmt"
	"time"
)

type Type string

const (
	ViewType Type = "ViewType"
	LikeType Type = "LikeType"
)

type Base struct {
	ID       string
	Type     Type
	DateTime time.Time
}

// Event ...
type Event interface {
	GetID() string
	GetType() Type
	SetID(id string)
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func New(t Type) (Event, error) {
	b := &Base{
		Type: t,
	}

	switch t {

	case ViewType:
		return &ViewEvent{
			Base: b,
		}, nil

	case LikeType:
		return &LikeEvent{
			Base: b,
		}, nil

	}

	return nil, fmt.Errorf("type %v not supported", t)
}

func (o *Base) GetID() string {
	return o.ID
}

func (o *Base) SetID(id string) {
	o.ID = id
}

func (o *Base) GetType() Type {
	return o.Type
}

func (o *Base) String() string {

	return fmt.Sprintf("id:%s type:%s", o.ID, o.Type)
}
