package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/port"
	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
)

type AlumnoUseCase struct {
	repo        port.AlumnoRepository
	fileStorage port.FileStorage
	notifier    port.NotificationService
}

func NewAlumnoUseCase(repo port.AlumnoRepository, fileStorage port.FileStorage, notifier port.NotificationService) *AlumnoUseCase {
	return &AlumnoUseCase{
		repo:        repo,
		fileStorage: fileStorage,
		notifier:    notifier,
	}
}

func (u *AlumnoUseCase) GetAll(ctx context.Context) ([]domain.Alumno, error) {
	return u.repo.GetAll(ctx)
}

func (u *AlumnoUseCase) GetByID(ctx context.Context, id uint) (*domain.Alumno, error) {
	alumno, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if alumno == nil {
		return nil, apperrors.ErrNotFound
	}
	return alumno, nil
}

func (u *AlumnoUseCase) Create(ctx context.Context, alumno *domain.Alumno) error {
	validationErrors := utils.ValidateAlumno(
		alumno.Nombres,
		alumno.Apellidos,
		alumno.Matricula,
		alumno.Promedio,
		alumno.Password,
		true,
	)
	if validationErrors.HasErrors() {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, validationErrors.Errors)
	}

	hashedPassword, err := utils.HashPassword(alumno.Password)
	if err != nil {
		return fmt.Errorf("error al hashear password: %w", err)
	}
	alumno.Password = hashedPassword

	return u.repo.Create(ctx, alumno)
}

func (u *AlumnoUseCase) Update(ctx context.Context, id uint, alumno *domain.Alumno) error {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.ErrNotFound
	}

	validationErrors := utils.ValidateAlumno(
		alumno.Nombres,
		alumno.Apellidos,
		alumno.Matricula,
		alumno.Promedio,
		alumno.Password,
		false,
	)
	if validationErrors.HasErrors() {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, validationErrors.Errors)
	}

	existing.Nombres = alumno.Nombres
	existing.Apellidos = alumno.Apellidos
	existing.Matricula = alumno.Matricula
	existing.Promedio = alumno.Promedio

	if alumno.Password != "" {
		hashedPassword, err := utils.HashPassword(alumno.Password)
		if err != nil {
			return fmt.Errorf("error al hashear password: %w", err)
		}
		existing.Password = hashedPassword
	}

	return u.repo.Update(ctx, existing)
}

func (u *AlumnoUseCase) Delete(ctx context.Context, id uint) error {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.ErrNotFound
	}

	return u.repo.Delete(ctx, id)
}

func (u *AlumnoUseCase) UploadFotoPerfil(ctx context.Context, id uint, file io.Reader, filename string, contentType string) (string, error) {
	alumno, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	if alumno == nil {
		return "", apperrors.ErrNotFound
	}

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("alumnos/%d/foto_perfil_%d%s", id, time.Now().Unix(), ext)

	url, err := u.fileStorage.Upload(ctx, key, file, contentType)
	if err != nil {
		return "", fmt.Errorf("error al subir foto: %w", err)
	}

	alumno.FotoPerfilUrl = url
	if err := u.repo.Update(ctx, alumno); err != nil {
		return "", fmt.Errorf("error al actualizar alumno: %w", err)
	}

	return url, nil
}

func (u *AlumnoUseCase) SendEmail(ctx context.Context, id uint) error {
	alumno, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if alumno == nil {
		return apperrors.ErrNotFound
	}

	subject := fmt.Sprintf("Calificaciones de %s %s", alumno.Nombres, alumno.Apellidos)
	message := fmt.Sprintf(
		"Información del alumno:\n\nNombre: %s %s\nMatrícula: %s\nPromedio: %.2f",
		alumno.Nombres,
		alumno.Apellidos,
		alumno.Matricula,
		alumno.Promedio,
	)

	if u.notifier == nil {
		return fmt.Errorf("notificador no configurado")
	}

	return u.notifier.Publish(ctx, subject, message)
}
