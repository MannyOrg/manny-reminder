package events

import (
	"github.com/gorilla/mux"
	"manny-reminder/pkg/models"
	"manny-reminder/pkg/utils"
	"net/http"
	"strconv"
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

func (h Handler) GetUsersEvents(w http.ResponseWriter, r *http.Request) {
	pt, s, err := h.getPagingData(r)
	if err != nil {
		utils.SendHttpError(w, err)
	}

	events, err := h.es.GetUsersEvents(pt, s)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	utils.SendJson(w, events)
}

func (h Handler) GetUserEvents(w http.ResponseWriter, r *http.Request) {
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

func (h Handler) getPagingData(r *http.Request) (string, int, error) {
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
