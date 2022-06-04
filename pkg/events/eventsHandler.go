package events

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"manny-reminder/pkg/models"
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

type GetUsersEventsResponse map[string][]models.Event

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

func (h Handler) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		utils.SendHttpStringError(w, "User id not defined")
		return
	}

	events, err := h.es.GetUserEvents(userId)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		utils.SendHttpError(w, err)
	}
}
