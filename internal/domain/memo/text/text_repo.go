package text

import (
	"context"
	"fmt"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/pkg/database/cockroach"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	"github.com/jackc/pgx/v5"
)

type TextRepo struct {
	etcd etcd.Etcd
	db   cockroach.Cockroach
}

func NewTextRepo(etcd etcd.Etcd, db cockroach.Cockroach) TextRepo {
	return TextRepo{
		etcd,
		db,
	}
}

func (r TextRepo) Save(ctx context.Context, record *memo.Memo) (*memo.Memo, error) {
	isNew := len(record.ID) == 0

	if isNew {
		id, err := memo.GenerateID()
		if err != nil {
			return nil, err
		}
		record = memo.NewMemo(record, memo.WithID(id))
	}
	err := r.db.ExecuteTx(ctx, func(tx pgx.Tx) error {
		if isNew {
			return r.insert(ctx, tx, record)
		} else {
			return r.update(ctx, tx, record)
		}
	})

	if err != nil {
		return nil, err
	}

	return record, nil
}

func (r TextRepo) insert(ctx context.Context, tx pgx.Tx, record *memo.Memo) error {
	_, err := tx.Exec(ctx, "INSERT INTO text (id, hash_code, kind, status, content) VALUES ($1, $2, $3, $4, $5)", record.ID, record.Hash, record.Kind, record.Status, record.Content)
	return err
}

func (r TextRepo) update(ctx context.Context, tx pgx.Tx, record *memo.Memo) error {
	_, err := tx.Exec(ctx, "UPDATE text SET status = $1 WHERE id = $2", record.Status, record.ID)
	return err
}

func (r TextRepo) FetchByID(ctx context.Context, id string) (*memo.Memo, error) {
	pool := r.db.Acquire()

	record := new(memo.Memo)

	if err := pool.QueryRow(ctx, "SELECT id, kind, status, content, hash_code FROM text WHERE id = $1", id).Scan(&record.ID, &record.Kind, &record.Status, &record.Content, &record.Hash); err != nil {
		fmt.Printf("Err %+v\n\n", err)
		return nil, &memo.NotFoundError{}
	}

	return record, nil
}

func (r TextRepo) FetchByHash(ctx context.Context, record *memo.Memo) (*memo.Memo, error) {
	pool := r.db.Acquire()

	m := new(memo.Memo)

	if err := pool.QueryRow(ctx, "SELECT id, kind, status, content, hash_code FROM text WHERE hash_code = $1", memo.Hash).Scan(&m.ID, &m.Kind, &m.Status, &m.Content, &m.Hash); err != nil {
		return nil, &memo.NotFoundError{}
	}

	return m, nil
}
