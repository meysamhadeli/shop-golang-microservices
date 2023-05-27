package rabbitmqcontainer

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RabbitMQ_Container(t *testing.T) {
	rabbitmqConn, _, _, err := Start(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, rabbitmqConn)
}
