package controller

import (
	"github.com/Bubotka/Microservices/proxy/internal/infrastructure/responder"
	"github.com/Bubotka/Microservices/proxy/internal/modules/user/service"
	"github.com/go-chi/chi"

	"net/http"
)

type Userer interface {
	GetByUsername(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	service service.Userer
	responder.Responder
}

func NewUserController(service service.Userer, responder responder.Responder) *UserController {
	return &UserController{service: service, Responder: responder}
}

func (u *UserController) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	out := u.service.GetByUsername(service.GetIn{Username: username})
	if out.Error != nil {
		u.Responder.ErrorBadRequest(w, out.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.Responder.OutputJSON(w, out.User)
}

func (u *UserController) List(w http.ResponseWriter, r *http.Request) {
	out := u.service.List()

	if out.Error != nil {
		u.Responder.ErrorBadRequest(w, out.Error)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.Responder.OutputJSON(w, ListResponse{
		Users: out.User,
	})
}
