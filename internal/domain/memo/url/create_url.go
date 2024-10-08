package url

import (
	"context"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type SaveUrlRepo interface {
	FetchByHash(context.Context, *memo.Memo) (*memo.Memo, error)
	Save(context.Context, *memo.Memo) (*memo.Memo, error)
}

type CreateUrlService struct {
	repo SaveUrlRepo
}

func NewCreateUrlService(repo SaveUrlRepo) CreateUrlService {
	return CreateUrlService{
		repo: repo,
	}
}

func (s CreateUrlService) GetKind() string {
	return "url"
}

func (s CreateUrlService) Save(ctx context.Context, record *memo.Memo) (*memo.Memo, error) {
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
