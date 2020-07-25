package handler

import (
	"errors"
	"fmt"

	"github.com/felipeagger/go-redis-streams/packages/event"
)

type viewHandler struct {
}

//NewViewHandler ...
func NewViewHandler() Handler {
	return &viewHandler{}
}

func (h *viewHandler) Handle(e event.Event) error {
	view, ok := e.(*event.ViewEvent)

	if !ok {
		return fmt.Errorf("incorrect event type")
	}

	if view.UserID == 5 {
		return errors.New("Falhou")
	}

	fmt.Printf("completed view %+v UserID: %v Extra:%v \n", view, view.UserID, view.Extra)

	return nil
}
