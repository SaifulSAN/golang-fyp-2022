package router

import (
	"database/sql"

	"decision-support-system-go/app/service/user"

	"github.com/go-chi/chi/v5"
)

// NewServer creates NewServer
func NewServer(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/register", func(r chi.Router) {
		r.Method("POST", "/user", user.RegisterUserStandard(db))
		r.Method("POST", "/organization", user.RegisterUserOrganization(db))
	})

	r.Route("/login", func(r chi.Router) {
		r.Method("POST", "/", user.LoginUser(db))
	})

	r.Route("/users", func(r chi.Router) {

		// /users GET, POST (get all user list, insert new user)
		// /users/{id} GET, PUT (get details, update details)
		// /users/{id}/password GET, PUT (get password for comparison, update password)

		//r.Use(middlewaretest.ExampleMiddleware)
		r.Method("GET", "/", user.RetrieveStuff(db))
		r.Method("POST", "/", user.PostStuff(db))
		r.Method("GET", "/{id}", user.RetrieveStuffSingular(db))
		r.Method("PUT", "/{id}", user.PutStuff(db))
	})

	return r
}
