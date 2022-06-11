package events

import (
	"github.com/gorilla/mux"
	"manny-reminder/internal/models"
	"manny-reminder/internal/utils"
	"net/http"
	"strconv"
)

type EventsHandler interface {
}

type HandlerImpl struct {
	es EventsService
}

func NewHandler(es EventsService) *HandlerImpl {
	return &HandlerImpl{es: es}
}

type GetUsersEventsResponse map[string][]models.Event

func (h HandlerImpl) GetUsersEvents(w http.ResponseWriter, r *http.Request) {
	_, s, err := h.getPagingData(r)
	if err != nil {
		utils.SendHttpError(w, err)
	}

	events, err := h.es.GetUsersEvents("", s)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	utils.SendJson(w, events)
}

func (h HandlerImpl) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["userId"]
	if userId == "" {
		utils.SendHttpStringError(w, "User id not defined")
		return
	}
	pt, s, err := h.getPagingData(r)
	if err != nil {
		utils.SendHttpError(w, err)
	}

	events, err := h.es.GetUserEvents(userId, pt, s)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}
	utils.SendJson(w, events)
}

func (h HandlerImpl) getPagingData(r *http.Request) (string, int, error) {
	size := 10
	var err error

	pageToken := r.URL.Query().Get("pageToken")

	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return "", 0, err
		}
	}

	return pageToken, size, nil
}
