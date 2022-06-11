package auth

import (
	"manny-reminder/internal/utils"
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
	utils.SendJson(w, users)
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	authUrl := h.as.GetTokenFromWeb()

	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func (h *Handler) SaveUser(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	err := h.as.SaveUser(authCode)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
