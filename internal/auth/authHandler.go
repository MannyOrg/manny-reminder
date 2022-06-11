package auth

import (
	"manny-reminder/internal/utils"
	"net/http"
)

type Handler interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type HandlerImpl struct {
	as AuthService
}

func NewHandler(as AuthService) *HandlerImpl {
	return &HandlerImpl{as: as}
}

func (h *HandlerImpl) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.as.GetUsers()
	if err != nil {
		utils.SendHttpError(w, err)
	}
	utils.SendJson(w, users)
}

func (h *HandlerImpl) AddUser(w http.ResponseWriter, r *http.Request) {
	authUrl := h.as.GetTokenFromWeb()

	http.Redirect(w, r, authUrl, http.StatusSeeOther)
}

func (h *HandlerImpl) SaveUser(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	err := h.as.SaveUser(authCode)
	if err != nil {
		utils.SendHttpError(w, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
