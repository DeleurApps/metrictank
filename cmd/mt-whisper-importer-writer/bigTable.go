package main

import (
	"github.com/grafana/metrictank/idx/bigtable"
	bigTableStore "github.com/grafana/metrictank/store/bigtable"
)

func newBigTableIndex(args []string) *bigtable.BigtableIdx {
	btFlags := bigtable.ConfigSetup()
	btFlags.Parse(args)
	bigtable.CliConfig.Enabled = true
	bigtable.ConfigProcess()

	return bigtable.New(bigtable.CliConfig)
}

func newBigTableStore(ttls []uint32, args []string) (*bigTableStore.Store, error) {
	storeConfigFlags := bigTableStore.ConfigSetup()
	storeConfigFlags.Parse(args)

	store, err := bigTableStore.NewStore(
		bigTableStore.CliConfig,
		ttls,
		uint32(bigTableStore.CliConfig.MaxChunkSpan.Seconds()))

	if err != nil {
		return nil, err
	}

	return store, nil
}
