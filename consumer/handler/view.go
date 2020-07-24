package handler

import (
	"fmt"

	"github.com/felipeagger/go-redis-streams/packages/event"
)

type viewHandler struct {
	state *State
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

	/*
		u, exist := h.state.Users[view.UserID]
		if !exist { // should have an event to create user before use
			u = &user.User{
				UseID:   view.UserID,
				Balance: 0,
			}
			h.state.Users[view.UserID] = u
		}

		u.Balance += view.Amount
	*/

	fmt.Printf("completed view %+v UserID: %v Extra:%v \n", view, view.UserID, view.Extra)

	return nil
}
