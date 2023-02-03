package events

import uuid "github.com/satori/go.uuid"

type ProductUpdated struct {
	ProductId uuid.UUID
}
