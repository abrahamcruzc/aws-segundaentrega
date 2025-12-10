package http

import (
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/http/handler"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	alumnoHandler   *handler.AlumnoHandler
	profesorHandler *handler.ProfesorHandler
	sesionHandler   *handler.SesionHandler
}

func NewRouter(
	alumnoHandler *handler.AlumnoHandler,
	profesorHandler *handler.ProfesorHandler,
	sesionHandler *handler.SesionHandler,
) *Router {
	return &Router{
		alumnoHandler:   alumnoHandler,
		profesorHandler: profesorHandler,
		sesionHandler:   sesionHandler,
	}
}

func (rt *Router) Setup() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares globales
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.CORS)
	r.Use(middleware.ContentType)

	// Rutas de alumnos
	r.Route("/alumnos", func(r chi.Router) {
		r.Get("/", rt.alumnoHandler.GetAll)
		r.Post("/", rt.alumnoHandler.Create)
		r.Get("/{id}", rt.alumnoHandler.GetByID)
		r.Put("/{id}", rt.alumnoHandler.Update)
		r.Delete("/{id}", rt.alumnoHandler.Delete)
		r.Post("/{id}/fotoPerfil", rt.alumnoHandler.UploadFotoPerfil)
		r.Post("/{id}/email", rt.alumnoHandler.SendEmail)

		// Rutas de sesión
		r.Post("/{id}/session/login", rt.sesionHandler.Login)
		r.Post("/{id}/session/verify", rt.sesionHandler.Verify)
		r.Post("/{id}/session/logout", rt.sesionHandler.Logout)
	})

	// Rutas de profesores
	r.Route("/profesores", func(r chi.Router) {
		r.Get("/", rt.profesorHandler.GetAll)
		r.Post("/", rt.profesorHandler.Create)
		r.Get("/{id}", rt.profesorHandler.GetByID)
		r.Put("/{id}", rt.profesorHandler.Update)
		r.Delete("/{id}", rt.profesorHandler.Delete)
	})

	return r
}
