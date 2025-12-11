package domain

import "time"

type Profesor struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	NumeroEmpleado int       `json:"numeroEmpleado" gorm:"not null;unique"`
	Nombres        string    `json:"nombres" gorm:"not null"`
	Apellidos      string    `json:"apellidos" gorm:"not null"`
	HorasClase     int       `json:"horasClase" gorm:"not null"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (Profesor) TableName() string {
	return "profesores"
}
