package server

import (
	"GoServer/api/handler"
	"GoServer/api/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run() error {
	r := chi.NewRouter()

	r.Use(middleware.CsrfMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", handler.SignInHandler)
			r.Post("/signup", handler.SignUpHandler)
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Get("/", handler.GetAllProfiles)
			r.Get("/{id}", handler.GetProfile)
			r.Route("/friends", func(r chi.Router) {
				r.Post("/", handler.AddFriend)
				r.Get("/", handler.GetAllFriends)
			})
		})

		r.Route("/conversations", func(r chi.Router) {
			r.Get("/", handler.GetAllConversations)
			r.Post("/", handler.CreateConversation)
			r.Route("/{id}", func(r chi.Router) {
				r.Route("/messages", func(r chi.Router) {
					r.Get("/", handler.GetAllMessages)
					r.Post("/", handler.CreateMessage)
					r.Get("/listener", handler.MessageListener)
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", handler.GetAllConversationUsers)
				})
			})
		})

		r.Get("/csrf", handler.GetCSRFToken)
	})

	r.Get("/*", handler.FrontendHandler)

	return http.ListenAndServe(":8080", r)
}
