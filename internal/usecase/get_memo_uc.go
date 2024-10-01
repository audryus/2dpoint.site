package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetMemo(id string) (Memo, error) {
	var m Memo
	endpoints := []string{"http://localhost:2379"}
	dialTimeout := 1 * time.Second

	// Create a new etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close() // Make sure to close the client

	// Set a new key-value pair
	ctx := context.Background()

	// Get the value for the key
	resp, err := cli.Get(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) > 0 {
		json.Unmarshal(resp.Kvs[0].Value, &m)

		put(cli, &m)
	} else {
		fmt.Printf("%s chave não encontrada ...\n", id) // Output: "bar"
	}

	return m, nil
}
