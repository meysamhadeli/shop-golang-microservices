package events

import uuid "github.com/satori/go.uuid"

type ProductCreated struct {
	ProductId uuid.UUID
}
