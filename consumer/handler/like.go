package handler

import (
	"errors"
	"fmt"

	evt "github.com/felipeagger/go-redis-streams/packages/event"
)

type likeHandler struct {
}

//NewLikeHandler ...
func NewLikeHandler() Handler {
	return &likeHandler{}
}

func (h *likeHandler) Handle(e evt.Event, retry bool) error {
	event, ok := e.(*evt.LikeEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	if event.UserID == 5 && !retry {
		return errors.New("Falhou")
	}

	fmt.Printf("completed like %+v UserID: %v\n", event, event.UserID)

	return nil
}
