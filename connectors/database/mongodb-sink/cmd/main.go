package main

import (
	"context"
	"github.com/linkall-labs/connector/mongodb-sink/internal"
	"os"
)

func main() {
	ga := internal.NewMongoSink(internal.Config{Port: 8081})
	err := ga.StartReceive(context.Background())
	if err != nil {
		//log.Error(context.Background(), "start controller proxy failed", map[string]interface{}{
		//	log.KeyError: err,
		//})
		os.Exit(-1)
	}

	err = ga.StartReceive(context.Background())
	if err != nil {
		//log.Error(context.Background(), "start CloudEvents gateway failed", map[string]interface{}{
		//	log.KeyError: err,
		//})
		os.Exit(-1)
	}
}
