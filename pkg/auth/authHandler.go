package auth

import (
	"encoding/json"
	"manny-reminder/pkg/utils"
	"net/http"
)

type IHandler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	as IService
}

func NewHandler(as IService) *Handler {
	return &Handler{as: as}
}

func (h *Handler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.as.GetUsers()
	if err != nil {
		utils.SendHttpError(w, err)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		utils.SendHttpError(w, err)
	}
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	authUrl := h.as.GetTokenFromWeb()

	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func (h *Handler) SaveUser(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	err := h.as.AddUser(authCode)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}
}