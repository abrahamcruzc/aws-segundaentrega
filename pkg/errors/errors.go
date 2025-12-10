package errors

import "errors"

var (
	ErrNotFound       = errors.New("recurso no encontrado")
	ErrInvalidInput   = errors.New("datos de entrada inválidos")
	ErrUnauthorized   = errors.New("no autorizado")
	ErrAlreadyExists  = errors.New("el recurso ya existe")
	ErrInternalServer = errors.New("error interno del servidor")
	ErrInvalidSession = errors.New("sesión inválida o expirada")
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v *ValidationErrors) Add(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}
