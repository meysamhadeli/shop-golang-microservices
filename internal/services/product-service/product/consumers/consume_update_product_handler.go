package consumers

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func HandleConsumeUpdateProduct(queue string, msg amqp.Delivery) error {

	log.Infof("Message received on queue: '%s' with message: %s", queue, string(msg.Body))
	return nil
}
