package handler

import (
	"errors"
	"fmt"

	"github.com/felipeagger/go-redis-streams/packages/event"
)

type likeHandler struct {
}

//NewLikeHandler ...
func NewLikeHandler() Handler {
	return &likeHandler{}
}

func (h *likeHandler) Handle(e event.Event) error {
	like, ok := e.(*event.LikeEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	if like.UserID == 5 {
		return errors.New("Falhou")
	}

	fmt.Printf("completed like %+v UserID: %v Extra:%v \n", like, like.UserID, like.Extra)

	return nil
}
