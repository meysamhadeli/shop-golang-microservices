package events

import uuid "github.com/satori/go.uuid"

type ProductDeleted struct {
	ProductId uuid.UUID
}
