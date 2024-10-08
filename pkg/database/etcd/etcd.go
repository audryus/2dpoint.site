package etcd

import (
	"context"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const dialTimeout = 1 * time.Second
const dialKeepAliveTime = 60 * time.Second

type Etcd interface {
	Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error)
	Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) error
	KeepAlive(ctx context.Context, id clientv3.LeaseID)
	Close() error
}

type etcd struct {
	cli *clientv3.Client
}

func New(cfg config.Config, l logger.Log) (Etcd, error) {
	addr := cfg.Etcd.Host + ":" + cfg.Etcd.Port

	endpoints := []string{addr}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		return nil, err
	}

	l.Info("etcd client created")

	return etcd{
		cli,
	}, nil
}

func (c etcd) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	return c.cli.Grant(ctx, ttl)
}
func (c etcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return c.cli.Get(ctx, key, opts...)
}

func (c etcd) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
	_, err := c.cli.Put(ctx, key, val, opts...)
	return err
}

func (c etcd) KeepAlive(ctx context.Context, id clientv3.LeaseID) {
	c.cli.KeepAlive(ctx, id)
}
func (c etcd) Close() error {
	return c.cli.Close()
}
