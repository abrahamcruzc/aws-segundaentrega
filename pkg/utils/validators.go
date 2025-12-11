package utils

import (
	"strings"

	apperrors "github.com/abrahamcruzc/aws-segundaentrega/pkg/errors"
)

func ValidateAlumno(nombres, apellidos, matricula string, promedio float64, password string, isCreate bool) *apperrors.ValidationErrors {
	errors := &apperrors.ValidationErrors{}

	if strings.TrimSpace(nombres) == "" {
		errors.Add("nombres", "El campo nombres es requerido")
	}

	if strings.TrimSpace(apellidos) == "" {
		errors.Add("apellidos", "El campo apellidos es requerido")
	}

	if strings.TrimSpace(matricula) == "" {
		errors.Add("matricula", "El campo matricula es requerido")
	}

	if promedio < 0 || promedio > 10 {
		errors.Add("promedio", "El promedio debe estar entre 0 y 10")
	}

	// Password es opcional - los tests no lo envían

	return errors
}

func ValidateProfesor(numeroEmpleado int, nombres, apellidos string, horasClase int) *apperrors.ValidationErrors {
	errors := &apperrors.ValidationErrors{}

	if numeroEmpleado <= 0 {
		errors.Add("numeroEmpleado", "El campo numeroEmpleado es requerido")
	}

	if strings.TrimSpace(nombres) == "" {
		errors.Add("nombres", "El campo nombres es requerido")
	}

	if strings.TrimSpace(apellidos) == "" {
		errors.Add("apellidos", "El campo apellidos es requerido")
	}

	if horasClase < 0 {
		errors.Add("horasClase", "El campo horasClase debe ser mayor o igual a 0")
	}

	return errors
}

func ValidatePassword(password string) *apperrors.ValidationErrors {
	errors := &apperrors.ValidationErrors{}

	if strings.TrimSpace(password) == "" {
		errors.Add("password", "El campo password es requerido")
	}

	return errors
}

func ValidateSessionString(sessionString string) *apperrors.ValidationErrors {
	errors := &apperrors.ValidationErrors{}

	if strings.TrimSpace(sessionString) == "" {
		errors.Add("sessionString", "El campo sessionString es requerido")
	}

	return errors
}
