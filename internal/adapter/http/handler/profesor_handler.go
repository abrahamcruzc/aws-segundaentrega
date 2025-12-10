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

type ProfesorHandler struct {
	service port.ProfesorService
}

func NewProfesorHandler(service port.ProfesorService) *ProfesorHandler {
	return &ProfesorHandler{service: service}
}

func (h *ProfesorHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profesores, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, profesores)
}

func (h *ProfesorHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	profesor, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Profesor no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, profesor)
}

func (h *ProfesorHandler) Create(w http.ResponseWriter, r *http.Request) {
	var profesor domain.Profesor
	if err := json.NewDecoder(r.Body).Decode(&profesor); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Create(r.Context(), &profesor); err != nil {
		if errors.Is(err, apperrors.ErrInvalidInput) {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, profesor)
}

func (h *ProfesorHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var profesor domain.Profesor
	if err := json.NewDecoder(r.Body).Decode(&profesor); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := h.service.Update(r.Context(), uint(id), &profesor); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Profesor no encontrado")
			return
		}
		if errors.Is(err, apperrors.ErrInvalidInput) {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Profesor actualizado correctamente")
}

func (h *ProfesorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Profesor no encontrado")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Profesor eliminado correctamente")
}
