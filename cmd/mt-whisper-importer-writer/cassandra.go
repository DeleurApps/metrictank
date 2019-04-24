package main

import (
	"fmt"

	"github.com/grafana/metrictank/idx/cassandra"
	cassandraStore "github.com/grafana/metrictank/store/cassandra"
)

func newCassandraIndex(args []string) *cassandra.CasIdx {
	cassFlags := cassandra.ConfigSetup()
	cassFlags.Parse(args)
	cassandra.CliConfig.Enabled = true
	cassandra.ConfigProcess()

	return cassandra.New(cassandra.CliConfig)
}

func newCassandraStore(ttls []uint32, args []string) (*cassandraStore.CassandraStore, error) {
	storeConfigFlags := cassandraStore.ConfigSetup()
	storeConfigFlags.Parse(args)

	// we don't use the cassandraStore's writeQueue, so we hard code this to 0
	cassandraStore.CliConfig.WriteQueueSize = 0

	store, err := cassandraStore.NewCassandraStore(cassandraStore.CliConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize cassandra store: %s", err)
	}

	return store, nil
}
