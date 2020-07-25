package handler

import (
	"fmt"

	"github.com/felipeagger/go-redis-streams/packages/event"
	evt "github.com/felipeagger/go-redis-streams/packages/event"
)

//HandlerFactory ...
func HandlerFactory() func(t event.Type) Handler {

	return func(t event.Type) Handler {
		switch t {
		case event.LikeType:
			return NewLikeHandler()
		case event.CommentType:
			return NewCommentHandler()
		default:
			return NewDefaultHandler()
		}
	}
}

type Handler interface {
	Handle(e event.Event, retry bool) error
}

type defaultHandler struct {
}

//NewViewHandler ...
func NewDefaultHandler() Handler {
	return &defaultHandler{}
}

func (h *defaultHandler) Handle(e evt.Event, retry bool) error {
	fmt.Printf("undefined event %+v\n", e)
	return nil
}
