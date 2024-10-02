package etcd

import (
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const dialTimeout = 1 * time.Second

func NewClient(cfg config.Config) (*clientv3.Client, error) {
	endpoints := []string{"http://localhost:2379"}

	// Create a new etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	return cli, nil
}
