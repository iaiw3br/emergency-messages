package router

import (
	"github.com/go-chi/chi/v5"
	"projects/emergency-messages/internal/controllers"
)

const (
	api = "/api/"
	v1  = api + "v1"
)

type Router struct {
	router   *chi.Mux
	message  *controllers.Message
	receiver *controllers.Receiver
	template *controllers.Template
}

func New(router *chi.Mux, message *controllers.Message, receiver *controllers.Receiver, template *controllers.Template) Router {
	return Router{
		router:   router,
		message:  message,
		receiver: receiver,
		template: template,
	}
}

func (r Router) Load() {
	r.router.Route(v1, func(router chi.Router) {
		router.Route("/messages", func(router chi.Router) {
			router.Post("/", r.message.Send)
		})
		router.Route("/templates", func(router chi.Router) {
			router.Post("/", r.template.Create)

			router.Route("/{id}", func(router chi.Router) {
				router.Delete("/", r.template.Delete)
				router.Patch("/", r.template.Update)
			})
		})
		router.Route("/receivers", func(router chi.Router) {
			router.Get("/city/:city", r.receiver.GetByCity)
			router.Post("/upload", r.receiver.Upload)
		})
	})
}
