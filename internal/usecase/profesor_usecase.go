package usecase

import (
	"context"
	"fmt"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/port"
	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
)

type ProfesorUseCase struct {
	repo port.ProfesorRepository
}

func NewProfesorUseCase(repo port.ProfesorRepository) *ProfesorUseCase {
	return &ProfesorUseCase{repo: repo}
}

func (u *ProfesorUseCase) GetAll(ctx context.Context) ([]domain.Profesor, error) {
	return u.repo.GetAll(ctx)
}

func (u *ProfesorUseCase) GetByID(ctx context.Context, id uint) (*domain.Profesor, error) {
	profesor, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if profesor == nil {
		return nil, apperrors.ErrNotFound
	}
	return profesor, nil
}

func (u *ProfesorUseCase) Create(ctx context.Context, profesor *domain.Profesor) error {
	validationErrors := utils.ValidateProfesor(
		profesor.NumeroEmpleado,
		profesor.Nombres,
		profesor.Apellidos,
		profesor.HorasClase,
	)
	if validationErrors.HasErrors() {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, validationErrors.Errors)
	}

	return u.repo.Create(ctx, profesor)
}

func (u *ProfesorUseCase) Update(ctx context.Context, id uint, profesor *domain.Profesor) error {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.ErrNotFound
	}

	validationErrors := utils.ValidateProfesor(
		profesor.NumeroEmpleado,
		profesor.Nombres,
		profesor.Apellidos,
		profesor.HorasClase,
	)
	if validationErrors.HasErrors() {
		return fmt.Errorf("%w: %v", apperrors.ErrInvalidInput, validationErrors.Errors)
	}

	existing.NumeroEmpleado = profesor.NumeroEmpleado
	existing.Nombres = profesor.Nombres
	existing.Apellidos = profesor.Apellidos
	existing.HorasClase = profesor.HorasClase

	return u.repo.Update(ctx, existing)
}

func (u *ProfesorUseCase) Delete(ctx context.Context, id uint) error {
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.ErrNotFound
	}

	return u.repo.Delete(ctx, id)
}
