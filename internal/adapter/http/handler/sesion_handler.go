package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/port"
	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type SesionHandler struct {
	service port.SesionService
}

func NewSesionHandler(service port.SesionService) *SesionHandler {
	return &SesionHandler{service: service}
}

type LoginRequest struct {
	Password string `json:"password"`
}

type SessionRequest struct {
	SessionString string `json:"sessionString"`
}

func (h *SesionHandler) Login(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.Password == "" {
		utils.JSONError(w, http.StatusBadRequest, "Password requerido")
		return
	}

	sesion, err := h.service.Login(r.Context(), uint(id), req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Alumno no encontrado")
			return
		}
		if errors.Is(err, apperrors.ErrUnauthorized) {
			utils.JSONError(w, http.StatusBadRequest, "Password incorrecto")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{
		"sessionString": sesion.SessionString,
	})
}

func (h *SesionHandler) Verify(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.SessionString == "" {
		utils.JSONError(w, http.StatusBadRequest, "SessionString requerido")
		return
	}

	if err := h.service.Verify(r.Context(), uint(id), req.SessionString); err != nil {
		if errors.Is(err, apperrors.ErrUnauthorized) || errors.Is(err, apperrors.ErrInvalidInput) {
			utils.JSONError(w, http.StatusBadRequest, "Sesión inválida o inactiva")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Sesión válida")
}

func (h *SesionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var req SessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if req.SessionString == "" {
		utils.JSONError(w, http.StatusBadRequest, "SessionString requerido")
		return
	}

	if err := h.service.Logout(r.Context(), uint(id), req.SessionString); err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Sesión no encontrada")
			return
		}
		if errors.Is(err, apperrors.ErrUnauthorized) {
			utils.JSONError(w, http.StatusBadRequest, "Sesión no pertenece al alumno")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONMessage(w, http.StatusOK, "Sesión cerrada correctamente")
}
