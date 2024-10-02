package domain

import (
	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type createMemoRepoEtcd struct {
	etcd *clientv3.Client
}

func NewCreateMemoRepoEtcd(etcd *clientv3.Client) createMemoRepoEtcd {
	return createMemoRepoEtcd{etcd: etcd}
}

func (r createMemoRepoEtcd) Put(memo Memo) error {
	cli := r.etcd
	id := memo.ID
	url := memo.Url

	b, err := json.Marshal(memo)

	ctx, timeout := context.WithTimeout(context.Background(), timeout)
	defer timeout()

	sec, _ := time.ParseDuration("3h")

	resp, err := cli.Grant(ctx, int64(sec.Seconds()))
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, id, string(b), clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, url, id, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	return nil
}
