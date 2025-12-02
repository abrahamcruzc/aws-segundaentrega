package postgres

import (
	"context"
	"errors"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"gorm.io/gorm"
)

type ProfesorRepository struct {
	db *gorm.DB
}

func NewProfesorRepository(db *gorm.DB) *ProfesorRepository {
	return &ProfesorRepository{db: db}
}

func (r *ProfesorRepository) GetAll(ctx context.Context) ([]domain.Profesor, error) {
	var profesores []domain.Profesor
	if err := r.db.WithContext(ctx).Find(&profesores).Error; err != nil {
		return nil, err
	}
	return profesores, nil
}

func (r *ProfesorRepository) GetByID(ctx context.Context, id uint) (*domain.Profesor, error) {
	var profesor domain.Profesor
	if err := r.db.WithContext(ctx).First(&profesor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profesor, nil
}

func (r *ProfesorRepository) Create(ctx context.Context, profesor *domain.Profesor) error {
	return r.db.WithContext(ctx).Create(profesor).Error
}

func (r *ProfesorRepository) Update(ctx context.Context, profesor *domain.Profesor) error {
	return r.db.WithContext(ctx).Save(profesor).Error
}

func (r *ProfesorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Profesor{}, id).Error
}
