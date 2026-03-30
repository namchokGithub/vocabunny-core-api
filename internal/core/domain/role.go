package domain

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID
	Code string
	Name string
	AuditFields
}
