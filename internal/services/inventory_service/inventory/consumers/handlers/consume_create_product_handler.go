package handlers

import (
	"encoding/json"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/consumers/events"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/inventory/models"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/inventory_service/shared/delivery"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func HandleConsumeCreateProduct(queue string, msg amqp.Delivery, inventoryDeliveryBase *delivery.InventoryDeliveryBase) error {

	log.Infof("Message received on queue: %s with message: %s", queue, string(msg.Body))

	var productCreated events.ProductCreated

	err := json.Unmarshal(msg.Body, &productCreated)
	if err != nil {
		return err
	}

	count := productCreated.Count

	productItem, _ := inventoryDeliveryBase.InventoryRepository.GetProductInInventories(inventoryDeliveryBase.Ctx, productCreated.ProductId)

	if productItem != nil {
		count = productItem.Count + count
	}

	_, err = inventoryDeliveryBase.InventoryRepository.AddProductItemToInventory(inventoryDeliveryBase.Ctx, &models.ProductItem{
		Id:          uuid.NewV4(),
		ProductId:   productCreated.ProductId,
		Count:       count,
		InventoryId: productCreated.InventoryId,
	})

	if err != nil {
		return err
	}

	return nil
}
