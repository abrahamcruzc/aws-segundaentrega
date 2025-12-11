package domain

type Sesion struct {
	ID            string `json:"id" dynamodbav:"id"`       // UUID
	Fecha         int64  `json:"fecha" dynamodbav:"fecha"` // UNIX TIMESTAMP
	AlumnoID      uint   `json:"alumnoId" dynamodbav:"alumnoId"`
	Active        bool   `json:"active" dynamodbav:"active"`
	SessionString string `json:"sessionString" dynamodbav:"sessionString"` // 128 CARACTERES ALEATORIOS
}
