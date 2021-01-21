package handlers

import (
	"net/http"

	getroutes "gitlab.myteksi.net/goscripts/zendesk/handlers/get-routes"
)

type IHandler interface {
	HandleGetRoutes(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	getRoutesHandler getroutes.IHandler
}

func NewHandlersImpl() IHandler {
	// Here the dependencies would be injected into the handler individually and then stored in Handlers struct
	getRouteHandler := getroutes.NewHandlerImpl()
	return &Handlers{getRoutesHandler: getRouteHandler}
}

func (h *Handlers) HandleGetRoutes(w http.ResponseWriter, r *http.Request) {
	h.getRoutesHandler.Handle(w, r)
}
