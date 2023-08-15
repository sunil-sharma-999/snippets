package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {

	mux := chi.NewRouter()

	// middlewares
	mux.Use(app.recoverPanic, app.logRequest, secureHeaders)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.session.Enable)
		mux.Use(noSurf)
		mux.Use(app.authenticate)
		// home
		mux.Get("/", app.home)
		// snippet
		mux.Get("/snippet/{id}", app.showSnippet)
		mux.Group(func(mux chi.Router) {
			mux.Use(app.requireAuthenticatedUser)
			mux.Get("/snippet/create", app.createSnippetForm)
			mux.Post("/snippet/create", app.createSnippet)
		})
		// auth
		mux.Get("/user/login", app.loginUserForm)
		mux.Post("/user/login", app.loginUser)
		mux.Get("/user/signup", app.signupUserForm)
		mux.Post("/user/signup", app.signupUser)
		mux.Post("/user/logout", app.logout)
	})

	mux.Get("/ping", ping)

	// static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
