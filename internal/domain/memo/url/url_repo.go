package url

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const HASH_KEY = "Hash:%s"

type urlRepo struct {
	etcd etcd.Etcd
}

func NewUrlRepo(etcd etcd.Etcd) urlRepo {
	return urlRepo{
		etcd,
	}
}

func (r urlRepo) Save(ctx context.Context, record *memo.Memo) (*memo.Memo, error) {
	cli := r.etcd

	isNew := len(record.ID) == 0

	if isNew {
		id, err := memo.GenerateID()
		if err != nil {
			return nil, err
		}
		record = memo.NewMemo(record, memo.WithID(id))
	}

	sec, _ := time.ParseDuration("3h")

	resp, err := cli.Grant(ctx, int64(sec.Seconds()))
	if err != nil {
		return nil, err
	}

	record = memo.NewMemo(record, memo.WithID(record.ID), memo.WithLease(resp.ID))

	b, err := json.Marshal(record)

	if err != nil {
		return nil, err
	}
	json := string(b)

	err = cli.Put(ctx, record.ID, json, clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}

	err = cli.Put(ctx, fmt.Sprintf(HASH_KEY, record.Hash), json, clientv3.WithLease(resp.ID))
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r urlRepo) getID(id string) (string, error) {
	if len(id) > 0 {
		return id, nil
	}
	return memo.GenerateID()
}

func (r urlRepo) FetchByID(ctx context.Context, id string) (*memo.Memo, error) {
	cli := r.etcd

	resp, err := cli.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.get(resp)
}

func (r urlRepo) get(resp *clientv3.GetResponse) (*memo.Memo, error) {
	if len(resp.Kvs) == 0 {
		return nil, &memo.NotFoundError{}
	}

	m := new(memo.Memo)

	if err := json.Unmarshal(resp.Kvs[0].Value, m); err != nil {
		fmt.Printf("wrong json format %s", string(resp.Kvs[0].Value))
		return nil, err
	}

	r.refreshTTL(m)
	return m, nil
}

func (r urlRepo) FetchByHash(ctx context.Context, u *memo.Memo) (*memo.Memo, error) {
	cli := r.etcd

	resp, err := cli.Get(ctx, fmt.Sprintf(HASH_KEY, u.Hash))
	if err != nil {
		return nil, err
	}

	return r.get(resp)
}

func (r urlRepo) refreshTTL(u *memo.Memo) {
	cli := r.etcd
	cli.KeepAlive(context.Background(), u.Lease)
}
