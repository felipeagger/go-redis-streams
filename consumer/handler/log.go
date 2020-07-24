package handler

import (
	"fmt"

	"github.com/felipeagger/go-redis-streams/packages/event"
)

type logHandler struct {
}

//NewLogHandler ...
func NewLogHandler() Handler {
	return &logHandler{}
}

func (h *logHandler) Handle(e event.Event) error {

	fmt.Printf("new event:%+v\n", e)
	return nil
}
