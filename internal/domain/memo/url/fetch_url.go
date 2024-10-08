package url

import (
	"context"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type GetUrlRepo interface {
	FetchByID(context.Context, string) (*memo.Memo, error)
}

type GetUrlService struct {
	repo GetUrlRepo
}

func NewGetUrlService(repo GetUrlRepo) GetUrlService {
	return GetUrlService{
		repo: repo,
	}
}

func (s GetUrlService) GetKind() string {
	return "url"
}

func (s GetUrlService) FetchByID(ctx context.Context, id string) (*memo.Memo, error) {
	record, err := s.repo.FetchByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return record, nil
}
