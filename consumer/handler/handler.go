package handler

import (
	"github.com/felipeagger/go-redis-streams/packages/event"
)

type State struct {
	LatestEventID string
}

func (s *State) SetLatestEventID(id string) {
	s.LatestEventID = id
}
func (s *State) GetLatestEventID() string {
	return s.LatestEventID
}

func HandlerFactory() func(t event.Type) Handler {

	return func(t event.Type) Handler {
		switch t {
		case event.ViewType:
			return NewViewHandler()
		case event.LikeType:
			return NewViewHandler()
		}
		return NewLogHandler()
	}
}

type Handler interface {
	Handle(e event.Event) error
}
