package events

import uuid "github.com/satori/go.uuid"

type InventoryUpdated struct {
	ProductId   uuid.UUID
	InventoryId int64
	Count       int32
}
