package handler

import (
	"Avito_Backend_Trainee/models"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"regexp"
)

var reportNameKey = "reportName"

func reports(router chi.Router) {
	router.Get("/", getUserReport)

	router.Route("/{reportName}", func(router chi.Router) {
		router.Use(ReportContext)
		router.Get("/", downloadReport)
	})
}

func ReportContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reportName := chi.URLParam(r, "reportName")
		if reportName == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("report name is required")))
			return
		}
		_, err := regexp.MatchString("[A-Za-z0-9]+_[A-Za-z0-9]+\\\\.csv", reportName)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("bad report name")))
		}
		ctx := context.WithValue(r.Context(), reportNameKey, reportName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func downloadReport(w http.ResponseWriter, r *http.Request) {
	reportName := r.Context().Value(reportNameKey).(string)
	contentType, filename, err := dbInstance.DownloadReportByName(reportName)
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	w.Header().Add("Content-Disposition", "attachment; filename="+reportName)
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, filename)
}

func getUserReport(w http.ResponseWriter, r *http.Request) {
	report := &models.UserReport{}
	if err := render.Bind(r, report); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := dbInstance.CreateReport(report); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, report); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
