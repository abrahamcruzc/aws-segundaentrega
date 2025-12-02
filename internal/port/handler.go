package port

import "net/http"

// AlumnoHandler - Endpoints HTTP para Alumno
type AlumnoHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	UploadFotoPerfil(w http.ResponseWriter, r *http.Request)
	SendEmail(w http.ResponseWriter, r *http.Request)
}

// ProfesorHandler - Endpoints HTTP para Profesor
type ProfesorHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// SessionHandler - Endpoints HTTP para Sesiones
type SesionHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}
