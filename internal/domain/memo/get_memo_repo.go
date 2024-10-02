package domain

import (
	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const timeout = 1 * time.Second

type getMemoRepoEtcd struct {
	etcd *clientv3.Client
}

func NewGetMemoRepoEtcd(etcd *clientv3.Client) getMemoRepoEtcd {
	return getMemoRepoEtcd{etcd: etcd}
}

func (r getMemoRepoEtcd) Get(id string) (Memo, error) {
	var m Memo
	cli := r.etcd

	ctx, timeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer timeout()

	// Get the value for the key
	resp, err := cli.Get(ctx, id)
	if err != nil {
		return m, err
	}

	if len(resp.Kvs) == 0 {
		return m, &NotFoundError{}
	}
	json.Unmarshal(resp.Kvs[0].Value, &m)
	return m, nil
}

func (r getMemoRepoEtcd) GetByUrl(url string) (Memo, error) {
	var m Memo
	cli := r.etcd

	ctx, timeout := context.WithTimeout(context.Background(), timeout)
	defer timeout()

	// Get the value for the key
	resp, err := cli.Get(ctx, url)
	if err != nil {
		return m, err
	}

	if len(resp.Kvs) == 0 {
		return m, &NotFoundError{}
	}

	id := string(resp.Kvs[0].Value)

	return r.Get(id)
}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "memo not found"
}
