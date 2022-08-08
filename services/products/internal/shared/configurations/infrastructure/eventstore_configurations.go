package infrastructure

import (
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/es/store"
	"github.com/meysamhadeli/shop-golang-microservices/pkg/eventstroredb"
)

func (ic *infrastructureConfigurator) configEventStore() (*esdb.Client, error, func()) {
	db, err := eventstroredb.NewEventStoreDB(ic.cfg.EventStoreConfig)
	if err != nil {
		return nil, err, nil
	}

	aggregateStore := store.NewAggregateStore(ic.log, db)
	fmt.Print(aggregateStore)

	return db, nil, func() {
		_ = db.Close() // nolint: errcheck
	}
}
