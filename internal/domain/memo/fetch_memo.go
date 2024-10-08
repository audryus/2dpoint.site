package memo

import (
	"context"
)

type FetchRepo interface {
	FetchKindByID(context.Context, string) (string, error)
}

type FetchKindRepo interface {
	GetKind() string
	FetchByID(context.Context, string) (*Memo, error)
}

type FetchService struct {
	repo  FetchRepo
	kinds map[string]FetchKindRepo
}

func NewFetchService(repo FetchRepo, kids ...FetchKindRepo) FetchService {
	kinds := make(map[string]FetchKindRepo)

	for _, k := range kids {
		kinds[k.GetKind()] = k
	}

	return FetchService{repo, kinds}
}

func (s FetchService) Fetch(ctx context.Context, id string) (*Memo, error) {
	kind, err := s.repo.FetchKindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.kinds[kind].FetchByID(ctx, id)
}
