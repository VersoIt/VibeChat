package http

import (
	"Messenger/internal/service"
	"net/http"
)

type Handler struct {
	authService *service.Auth
	server      *http.Server
	handler     *http.ServeMux
}

func NewHandler(authService *service.Auth) *Handler {
	h := &Handler{authService: authService, handler: http.NewServeMux()}
	h.registerHandlers()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP(w, r)
}

func (h *Handler) registerHandlers() {
	h.handler.HandleFunc("/users", h.getUserInfo)
	h.handler.HandleFunc("/sign-up", h.signUp)
	h.handler.HandleFunc("/sign-in-email", h.signInViaEmail)
	h.handler.HandleFunc("/sign-in-login", h.signInViaLogin)
	h.handler.HandleFunc("/sign-out", h.signOut)
	h.handler.HandleFunc("/validate", h.validateAccessToken)
	h.handler.HandleFunc("/refresh", h.refreshAccessToken)
}
