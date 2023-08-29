package handler

import (
	"Avito_Backend_Trainee/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

func segments(router chi.Router) {
	router.Route("/", func(router chi.Router) {
		router.Post("/", createSegment)
		router.Delete("/", deleteSegment)
	})
}

func createSegment(w http.ResponseWriter, r *http.Request) {
	segment := &models.Segment{}
	if err := render.Bind(r, segment); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.AddSegment(segment); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segment); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteSegment(w http.ResponseWriter, r *http.Request) {
	segment := &models.Segment{}
	if err := render.Bind(r, segment); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.DeleteSegment(segment); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, segment); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
