package text

import (
	"context"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type GetTextRepo interface {
	FetchByID(context.Context, string) (*memo.Memo, error)
}

type GetTextService struct {
	repo GetTextRepo
}

func NewGetTextService(repo GetTextRepo) GetTextService {
	return GetTextService{
		repo: repo,
	}
}

func (s GetTextService) GetKind() string {
	return "text"
}

func (s GetTextService) FetchByID(ctx context.Context, id string) (*memo.Memo, error) {
	memo, err := s.repo.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return memo, nil
}
