package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grafana/metrictank/mdata"
	bigTableStore "github.com/grafana/metrictank/store/bigtable"
	cassandraStore "github.com/grafana/metrictank/store/cassandra"

	"github.com/grafana/metrictank/cluster"
	"github.com/grafana/metrictank/cluster/partitioner"
	"github.com/grafana/metrictank/idx"
	"github.com/grafana/metrictank/idx/bigtable"
	"github.com/grafana/metrictank/idx/cassandra"
	"github.com/grafana/metrictank/logger"
	"github.com/raintank/dur"
	"github.com/raintank/schema"
	log "github.com/sirupsen/logrus"
)

type storeType uint8

const (
	cass storeType = iota + 1
	bt
)

var (
	globalFlags = flag.NewFlagSet("global config flags", flag.ExitOnError)

	exitOnError = globalFlags.Bool(
		"exit-on-error",
		true,
		"Exit with a message when there's an error",
	)
	verbose = globalFlags.Bool(
		"verbose",
		false,
		"More detailed logging",
	)
	httpEndpoint = globalFlags.String(
		"http-endpoint",
		"127.0.0.1:8080",
		"The http endpoint to listen on",
	)
	ttlsStr = globalFlags.String(
		"ttls",
		"35d",
		"list of ttl strings used by MT separated by ','",
	)
	partitionScheme = globalFlags.String(
		"partition-scheme",
		"bySeries",
		"method used for partitioning metrics. This should match the settings of tsdb-gw. (byOrg|bySeries)",
	)
	uriPath = globalFlags.String(
		"uri-path",
		"/chunks",
		"the URI on which we expect chunks to get posted",
	)
	numPartitions = globalFlags.Int(
		"num-partitions",
		1,
		"Number of Partitions",
	)

	version = "(none)"
)

type Server struct {
	partitioner partitioner.Partitioner
	store       mdata.Store
	index       idx.MetricIndex
	HTTPServer  *http.Server
}

func init() {
	formatter := &logger.TextFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05.000"
	log.SetFormatter(formatter)
	log.SetLevel(log.InfoLevel)
}

func main() {

	flag.Usage = func() {
		fmt.Println("mt-whisper-importer-writer")
		fmt.Println()
		fmt.Println("Opens an endpoint to send data to, which then gets stored in the MT internal DB(s)")
		fmt.Println()
		fmt.Printf("Usage:\n\n")
		fmt.Printf("  mt-whisper-importer-writer [global config flags] <storetype> [store config flags] <idxtype> [idx config flags]\n\n")
		fmt.Printf("global config flags:\n\n")
		globalFlags.PrintDefaults()
		fmt.Println()
		fmt.Printf("storetype: only 'cass' and 'bt' supported\n\n")
		fmt.Printf("cass store config flags:\n\n")
		cassandraStore.ConfigSetup().PrintDefaults()
		fmt.Println()
		fmt.Printf("bt store config flags:\n\n")
		bigTableStore.ConfigSetup().PrintDefaults()
		fmt.Println()
		fmt.Printf("idxtype: only 'cass' and 'bt' supported\n\n")
		fmt.Printf("cass index config flags:\n\n")
		cassandra.ConfigSetup().PrintDefaults()
		fmt.Println()
		fmt.Printf("bt index config flags:\n\n")
		bigtable.ConfigSetup().PrintDefaults()
		fmt.Println()
		fmt.Println("EXAMPLES:")
		fmt.Println("mt-whisper-importer-writer -exit-on-error=true -fake-avg-aggregates=true -http-endpoint=0.0.0.0:8080 -num-partitions=8 -partition-scheme=bySeries -ttls=8d,2y -uri-path=/chunks -verbose=true cass -addrs=192.168.0.1 -keyspace=mydata -window-factor=20 cass -hosts=192.168.0.1:9042 -keyspace=mydata")
	}

	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.Usage()
		os.Exit(0)
	}

	var storeI, indexI int
	var useStore, useIndex storeType
	for i, v := range os.Args {
		if v == "cass" || v == "bt" {
			if storeI == 0 {
				storeI = i
				if v == "cass" {
					useStore = cass
				} else {
					useStore = bt
				}
			} else if indexI == 0 {
				indexI = i
				if v == "cass" {
					useStore = cass
				} else {
					useStore = bt
				}
			} else {
				fmt.Println("store & index types can only be specified once for each")
				flag.Usage()
				os.Exit(1)
			}
		}
	}

	if indexI == 0 {
		fmt.Println("store & index types need to be specified, they can be \"cass\" or \"bt\"")
		flag.Usage()
		os.Exit(1)
	}

	globalFlags.Parse(os.Args[1:storeI])
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	splits := strings.Split(*ttlsStr, ",")
	ttls := make([]uint32, 0)
	for _, split := range splits {
		ttls = append(ttls, dur.MustParseNDuration("ttl", split))
	}

	var index idx.MetricIndex
	var err error
	switch useIndex {
	case cass:
		index = newCassandraIndex(os.Args[storeI+1 : indexI])
	case bt:
		index = newBigTableIndex(os.Args[storeI+1 : indexI])
	}

	var store mdata.Store
	switch useStore {
	case cass:
		store, err = newCassandraStore(ttls, os.Args[indexI+1:])
		if err != nil {
			log.Fatalf("Failed to initialize cassandra store: %s", err)
		}
	case bt:
		store, err = newBigTableStore(ttls, os.Args[indexI+1:])
		if err != nil {
			log.Fatalf("Failed to initialize the bigtable store: %s", err)
		}
	}

	p, err := partitioner.NewKafka(*partitionScheme)
	if err != nil {
		panic(fmt.Sprintf("Failed to instantiate partitioner: %q", err))
	}

	cluster.Init("mt-whisper-importer-writer", version, time.Now(), "http", int(80))

	server := &Server{
		partitioner: p,
		index:       index,
		HTTPServer: &http.Server{
			Addr:        *httpEndpoint,
			ReadTimeout: 10 * time.Minute,
		},
		store: store,
	}
	server.index.Init()

	http.HandleFunc(*uriPath, server.chunksHandler)
	http.HandleFunc("/healthz", server.healthzHandler)

	log.Infof("Listening on %q", *httpEndpoint)
	err = http.ListenAndServe(*httpEndpoint, nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating listener: %q", err))
	}
}

func throwError(msg string) {
	msg = fmt.Sprintf("%s\n", msg)
	if *exitOnError {
		log.Panic(msg)
	} else {
		log.Error(msg)
	}
}

func (s *Server) healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}

func (s *Server) chunksHandler(w http.ResponseWriter, req *http.Request) {
	data := mdata.ArchiveRequest{}
	err := data.UnmarshalCompressed(req.Body)
	if err != nil {
		throwError(fmt.Sprintf("Error decoding cwr stream: %q", err))
		return
	}

	if len(data.ChunkWriteRequests) == 0 {
		log.Warn("Received empty list of cwrs")
		return
	}

	log.Debugf(
		"Received %d cwrs for metric %s. The first has Key: %s, T0: %d, TTL: %d. The last has Key: %s, T0: %d, TTL: %d",
		len(data.ChunkWriteRequests),
		data.MetricData.Name,
		data.ChunkWriteRequests[0].Key.String(),
		data.ChunkWriteRequests[0].T0,
		data.ChunkWriteRequests[0].TTL,
		data.ChunkWriteRequests[len(data.ChunkWriteRequests)-1].Key.String(),
		data.ChunkWriteRequests[len(data.ChunkWriteRequests)-1].T0,
		data.ChunkWriteRequests[len(data.ChunkWriteRequests)-1].TTL)

	partition, err := s.partitioner.Partition(&data.MetricData, int32(*numPartitions))
	if err != nil {
		throwError(fmt.Sprintf("Error partitioning: %q", err))
		return
	}

	mkey, err := schema.MKeyFromString(data.MetricData.Id)
	if err != nil {
		throwError(fmt.Sprintf("Received invalid id: %s", data.MetricData.Id))
		return
	}

	s.index.AddOrUpdate(mkey, &data.MetricData, partition)
	for _, cwr := range data.ChunkWriteRequests {
		s.store.Add(&cwr)
	}
}
