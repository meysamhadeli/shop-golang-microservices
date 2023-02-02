package container

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RabbitMQ_Container(t *testing.T) {
	rabbitmqConn, err := NewRabbitMQTestContainers().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, rabbitmqConn)
}
