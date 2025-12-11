package domain

import "time"

type Alumno struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Nombres       string    `json:"nombres" gorm:"not null"`
	Apellidos     string    `json:"apellidos" gorm:"not null"`
	Matricula     string    `json:"matricula" gorm:"not null;unique"`
	Promedio      float64   `json:"promedio" gorm:"not null"`
	FotoPerfilUrl string    `json:"fotoPerfilUrl,omitempty"`
	Password      string    `json:"-" gorm:"not null"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}
