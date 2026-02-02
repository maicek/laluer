package handler

import (
	"fmt"

	"github.com/maicek/laluer/core/apps"
)

type Action struct {
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}

func (s *HandlerService) Call(action Action) {
	switch action.Event {
	case "run":
		app := apps.GetApplcationByPath(action.Payload.(map[string]interface{})["path"].(string))
		if app != nil {
			app.Run()
		}
	default:
		fmt.Printf("Unknown action: %s, payload: %+v\n", action.Event, action.Payload)
		return
	}
}
