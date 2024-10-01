package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/audryus/2dpoint.site/internal/app"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Set up etcd client configuration
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
	_, err = cli.Put(ctx, "foo", "https://example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Get the value for the key
	ctx = context.Background()
	resp, err := cli.Get(ctx, "fd7801fb")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp.Kvs[0].Value)) // Output: "bar"

	// Delete the key
	_, err = cli.Delete(ctx, "foo")
	if err != nil {
		log.Fatal(err)
	}

	/*fiber := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = fiber.Shutdown()
	}()

	if err := fiber.Listen(":3000"); err != nil {
		log.Panic(err)
	}*/

	app.Run()

}

type Memo struct {
	Type string `json:"type" form:"type"`
}
