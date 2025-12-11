package postgres

import (
	"context"
	"errors"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"gorm.io/gorm"
)

type AlumnoRepository struct {
	db *gorm.DB
}

func NewAlumnoRepository(db *gorm.DB) *AlumnoRepository {
	return &AlumnoRepository{db: db}
}

func (r *AlumnoRepository) GetAll(ctx context.Context) ([]domain.Alumno, error) {
	var alumnos []domain.Alumno
	if err := r.db.WithContext(ctx).Find(&alumnos).Error; err != nil {
		return nil, err
	}
	return alumnos, nil
}

func (r *AlumnoRepository) GetByID(ctx context.Context, id uint) (*domain.Alumno, error) {
	var alumno domain.Alumno
	if err := r.db.WithContext(ctx).First(&alumno, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &alumno, nil
}

func (r *AlumnoRepository) Create(ctx context.Context, alumno *domain.Alumno) error {
	return r.db.WithContext(ctx).Create(alumno).Error
}

func (r *AlumnoRepository) Update(ctx context.Context, alumno *domain.Alumno) error {
	return r.db.WithContext(ctx).Save(alumno).Error
}

func (r *AlumnoRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Alumno{}, id).Error
}
