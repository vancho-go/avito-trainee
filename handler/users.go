package handler

import (
	"Avito_Backend_Trainee/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

func users(router chi.Router) {
	router.Route("/", func(router chi.Router) {
		router.Get("/", getActiveSegments)
		router.Post("/", createUser)
		router.Put("/", updateUser)
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddUser(user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userSegments := &models.UserSegmentsUpdate{}
	if err := render.Bind(r, userSegments); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.UpdateUserSegments(userSegments); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, userSegments); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getActiveSegments(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	segms, err := dbInstance.GetUserActiveSegments(user)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segms); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}
