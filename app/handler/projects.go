package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mingrammer/go-todo-rest-api-example/app/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "GetAllProjects")
	defer span.End()

	projects := []model.Project{}
	db.Find(&projects)
	respondJSON(w, http.StatusOK, projects)
}

func CreateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "CreateProject")
	defer span.End()

	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, project)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "GetProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]

	span.SetAttributes(attribute.String("project.title", title))

	project := getProjectOr404(ctx, db, title, w, r)
	if project == nil {
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func UpdateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "UpdateProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]

	span.SetAttributes(attribute.String("project.title", title))

	project := getProjectOr404(ctx, db, title, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "DeleteProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]

	span.SetAttributes(attribute.String("project.title", title))

	project := getProjectOr404(ctx, db, title, w, r)
	if project == nil {
		return
	}
	if err := db.Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "ArchiveProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]

	span.SetAttributes(attribute.String("project.title", title))

	project := getProjectOr404(ctx, db, title, w, r)
	if project == nil {
		return
	}
	project.Archive()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("handler")
	ctx, span := tracer.Start(r.Context(), "RestoreProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]

	span.SetAttributes(attribute.String("project.title", title))

	project := getProjectOr404(ctx, db, title, w, r)
	if project == nil {
		return
	}
	project.Restore()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	