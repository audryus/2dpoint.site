package memo

import (
	"context"
	"fmt"

	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const KIND = "KIND:%s"

type memoRepo struct {
	etcd etcd.Etcd
}

func NewMemoRepo(etcd etcd.Etcd) memoRepo {
	return memoRepo{etcd}
}

func (r memoRepo) FetchKindByID(ctx context.Context, id string) (string, error) {
	cli := r.etcd

	resp, err := cli.Get(ctx, fmt.Sprintf(KIND, id))
	if err != nil {
		return "", &NotFoundError{}
	}

	if len(resp.Kvs) == 0 {
		return "", &NotFoundError{}
	}

	return string(resp.Kvs[0].Value), nil
}

func (r memoRepo) SaveKind(ctx context.Context, memo *Memo) error {
	cli := r.etcd

	if err := cli.Put(ctx, fmt.Sprintf(KIND, memo.ID), memo.Kind, clientv3.WithLease(memo.Lease)); err != nil {
		return err
	}

	return nil
}
