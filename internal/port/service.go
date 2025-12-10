package port

import (
	"context"
	"io"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
)

// AlumnoService - Lógica de negocio para Alumno
type AlumnoService interface {
	GetAll(ctx context.Context) ([]domain.Alumno, error)
	GetByID(ctx context.Context, id uint) (*domain.Alumno, error)
	Create(ctx context.Context, alumno *domain.Alumno) error
	Update(ctx context.Context, id uint, alumno *domain.Alumno) error
	Delete(ctx context.Context, id uint) error
	UploadFotoPerfil(ctx context.Context, id uint, file io.Reader, filename string, contentType string) (string, error)
	SendEmail(ctx context.Context, id uint) error
}

// ProfesorService - Lógica de negocio para Profesor
type ProfesorService interface {
	GetAll(ctx context.Context) ([]domain.Profesor, error)
	GetByID(ctx context.Context, id uint) (*domain.Profesor, error)
	Create(ctx context.Context, profesor *domain.Profesor) error
	Update(ctx context.Context, id uint, profesor *domain.Profesor) error
	Delete(ctx context.Context, id uint) error
}

type SesionService interface {
	Login(ctx context.Context, alumnoID uint, password string) (*domain.Sesion, error)
	Verify(ctx context.Context, alumnoID uint, sessionString string) error
	Logout(ctx context.Context, alumnoID uint, sessionString string) error
}
