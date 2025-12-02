package postgres

import (
	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Alumno{},
		&domain.Profesor{},
	)
}

