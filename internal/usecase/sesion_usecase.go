package usecase

import (
	"context"
	"time"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/port"
	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
	"github.com/abrahamcruzc/aws-segundaentrega/pkg/utils"
	"github.com/google/uuid"
)

type SesionUseCase struct {
	sesionRepo port.SesionRepository
	alumnoRepo port.AlumnoRepository
}

func NewSesionUseCase(sesionRepo port.SesionRepository, alumnoRepo port.AlumnoRepository) *SesionUseCase {
	return &SesionUseCase{
		sesionRepo: sesionRepo,
		alumnoRepo: alumnoRepo,
	}
}

func (u *SesionUseCase) Login(ctx context.Context, alumnoID uint, password string) (*domain.Sesion, error) {
	alumno, err := u.alumnoRepo.GetByID(ctx, alumnoID)
	if err != nil {
		return nil, err
	}
	if alumno == nil {
		return nil, apperrors.ErrNotFound
	}

	if !utils.CheckPassword(alumno.Password, password) {
		return nil, apperrors.ErrUnauthorized
	}

	sessionString, err := utils.GenerateSessionString()
	if err != nil {
		return nil, err
	}

	sesion := &domain.Sesion{
		ID:            uuid.New().String(),
		Fecha:         time.Now().Unix(),
		AlumnoID:      alumnoID,
		Active:        true,
		SessionString: sessionString,
	}

	if err := u.sesionRepo.Create(ctx, sesion); err != nil {
		return nil, err
	}

	return sesion, nil
}

func (u *SesionUseCase) Verify(ctx context.Context, alumnoID uint, sessionString string) error {
	validationErrors := utils.ValidateSessionString(sessionString)
	if validationErrors.HasErrors() {
		return apperrors.ErrInvalidInput
	}

	sesion, err := u.sesionRepo.GetBySessionString(ctx, sessionString)
	if err != nil {
		return err
	}
	if sesion == nil {
		return apperrors.ErrUnauthorized
	}

	if sesion.AlumnoID != alumnoID {
		return apperrors.ErrUnauthorized
	}

	if !sesion.Active {
		return apperrors.ErrUnauthorized
	}

	return nil
}

func (u *SesionUseCase) Logout(ctx context.Context, alumnoID uint, sessionString string) error {
	validationErrors := utils.ValidateSessionString(sessionString)
	if validationErrors.HasErrors() {
		return apperrors.ErrInvalidInput
	}

	sesion, err := u.sesionRepo.GetBySessionString(ctx, sessionString)
	if err != nil {
		return err
	}
	if sesion == nil {
		return apperrors.ErrNotFound
	}

	if sesion.AlumnoID != alumnoID {
		return apperrors.ErrUnauthorized
	}

	return u.sesionRepo.Deactivate(ctx, sesion.ID)
}
