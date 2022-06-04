package events

import (
	"encoding/json"
	"manny-reminder/pkg/utils"
	"net/http"
)

type IHandler interface {
}

type Handler struct {
	es IService
}

func NewHandler(es IService) *Handler {
	return &Handler{es: es}
}

func (h Handler) GetUsersEvents(w http.ResponseWriter, _ *http.Request) {
	events, err := h.es.GetUsersEvents()
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		utils.SendHttpError(w, err)
	}
}
