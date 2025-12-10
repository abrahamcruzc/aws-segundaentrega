package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/port"
	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type AlumnoHandler struct {
	service port.AlumnoService
}

func NewAlumnoHandler(service port.AlumnoService) *AlumnoHandler {
	return &AlumnoHandler{service: service}
}

func (h *AlumnoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	alumnos, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, alumnos)
}

func (h *AlumnoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	alumno, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, alumno)
}

func (h *AlumnoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var alumno domain.Alumno
	if err := json.NewDecoder(r.Body).Decode(&alumno); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &alumno); err != nil {
		if errors.Is(err, apperrors.ErrInvalidInput) {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, alumno)
}

func (h *AlumnoHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var alumno domain.Alumno
	if err := json.NewDecoder(r.Body).Decode(&alumno); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Update(r.Context(), uint(id), &alumno); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		if errors.Is(err, apperrors.ErrInvalidInput) {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Alumno actualizado correctamente")
}

func (h *AlumnoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Alumno eliminado correctamente")
}

func (h *AlumnoHandler) UploadFotoPerfil(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Error al procesar el archivo")
		return
	}

	file, header, err := r.FormFile("foto")
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Archivo 'foto' requerido")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	url, err := h.service.UploadFotoPerfil(r.Context(), uint(id), file, header.Filename, contentType)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{"fotoPerfilUrl": url})
}

func (h *AlumnoHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.SendEmail(r.Context(), uint(id)); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Email enviado correctamente")
}
