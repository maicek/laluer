package handler

import (
	"fmt"

	"github.com/maicek/laluer/core/apps"
)

type Action struct {
	Event   string
	Payload any
}

type ActionRunPayload struct {
	Path string
}

func (s *HandlerService) Call(action Action) {
	switch action.Event {
	case "run":
		payload := action.Payload.(ActionRunPayload)
		fmt.Printf("Running app: %s\n", payload.Path)
		app := apps.GetApplcationByPath(payload.Path)
		if app != nil {
			app.Run()
		}
	default:
		fmt.Printf("Unknown action: %s, payload: %+v\n", action.Event, action.Payload)
		return
	}
}
