package handler

import (
	"errors"
	"fmt"

	evt "github.com/felipeagger/go-redis-streams/packages/event"
)

type commentHandler struct {
}

//NewCommentHandler ...
func NewCommentHandler() Handler {
	return &commentHandler{}
}

func (h *commentHandler) Handle(e evt.Event, retry bool) error {
	event, ok := e.(*evt.CommentEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	if event.UserID == 5 && !retry {
		return errors.New("Falhou")
	}

	fmt.Printf("processed event %+v UserID: %v Comment:%v \n", event, event.UserID, event.Comment)

	return nil
}
