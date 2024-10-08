package text

import (
	"context"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type SaveTextRepo interface {
	FetchByHash(ctx context.Context, u *memo.Memo) (*memo.Memo, error)
	Save(context.Context, *memo.Memo) (*memo.Memo, error)
}

type CreateTextService struct {
	repo SaveTextRepo
}

func NewCreateTextService(repo SaveTextRepo) CreateTextService {
	return CreateTextService{
		repo: repo,
	}
}

func (s CreateTextService) GetKind() string {
	return "text"
}

func (s CreateTextService) Save(ctx context.Context, record *memo.Memo) (*memo.Memo, error) {
	oldMemo, err := s.repo.FetchByHash(ctx, record)

	if err != nil {
		_, ok := err.(*memo.NotFoundError)
		if !ok {
			return nil, err
		}
	} else {
		return oldMemo, nil
	}

	return s.repo.Save(ctx, record)
}
