package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mingrammer/go-todo-rest-api-example/app/model"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/mingrammer/go-todo-rest-api-example")

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetAllProjects")
	defer span.End()

	projects := []model.Project{}
	db.Find(&projects)
	respondJSON(w, http.StatusOK, projects)
}

func CreateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "CreateProject")
	defer span.End()

	project := model.Project{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.WithContext(ctx).Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, project)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "GetProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]
	project := getProjectOr404(db.WithContext(ctx), title, w, r)
	if project == nil {
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func UpdateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "UpdateProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]
	project := getProjectOr404(db.WithContext(ctx), title, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.WithContext(ctx).Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "DeleteProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]
	project := getProjectOr404(db.WithContext(ctx), title, w, r)
	if project == nil {
		return
	}
	if err := db.WithContext(ctx).Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "ArchiveProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]
	project := getProjectOr404(db.WithContext(ctx), title, w, r)
	if project == nil {
		return
	}
	project.Archive()
	if err := db.WithContext(ctx).Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "RestoreProject")
	defer span.End()

	vars := mux.Vars(r)
	title := vars["title"]
	project := getProjectOr404(db.WithContext(ctx), title, w, r)
	if project == nil {
		return
	}
	project.Restore()
	if err := db.WithContext(ctx).Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

// getProjectOr404 gets a project instance if exists, or respond the 404 error otherwise
func getProjectOr404(db *gorm.DB, title string, w http.ResponseWriter, r *http.Request) *model.Project {
	_, span := tracer.Start(r.Context(), "getProjectOr404")
	defer span.End()

