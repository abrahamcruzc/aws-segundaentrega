package port

import (
	"context"
	"io"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
)

// AlumnoRepository - Operaciones de persistencia para Alumno
type AlumnoRepository interface {
	GetAll(ctx context.Context) ([]domain.Alumno, error)
	GetByID(ctx context.Context, id uint) (*domain.Alumno, error)
	Create(ctx context.Context, alumno *domain.Alumno) error
	Update(ctx context.Context, alumno *domain.Alumno) error
	Delete(ctx context.Context, id uint) error
}

// ProfesorRepository - Operaciones de persistencia para Profesor
type ProfesorRepository interface {
	GetAll(ctx context.Context) ([]domain.Profesor, error)
	GetByID(ctx context.Context, id uint) (*domain.Profesor, error)
	Create(ctx context.Context, profesor *domain.Profesor) error
	Update(ctx context.Context, profesor *domain.Profesor) error
	Delete(ctx context.Context, id uint) error
}

// SesionRepository - Operaciones de persistencia para Sesión
type SesionRepository interface {
	Create(ctx context.Context, sesion *domain.Sesion) error
	GetBySessionString(ctx context.Context, sessionString string) (*domain.Sesion, error)
	Deactivate(ctx context.Context, sessionString string) error
}

// FileStorage - Operaciones de alamacenamiento de archivos
type FileStorage interface {
	Upload(ctx context.Context, key string, file io.Reader, contentType string) (string, error)
	GetURL(ctx context.Context, key string) string
	Delete(ctx context.Context, key string) error
}

// NotificationService - Operaciones de notificación
type NotificationService interface {
	SendEmail(ctx context.Context, email string, subject string, message string) error
}
