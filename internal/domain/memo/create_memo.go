package memo

import (
	"context"
)

type SaveRepo interface {
	SaveKind(context.Context, *Memo) error
	FetchKindByID(context.Context, string) (string, error)
}

type SaveKindService interface {
	GetKind() string
	Save(context.Context, *Memo) (*Memo, error)
}

type CreateMemoService struct {
	repo  SaveRepo
	kinds map[string]SaveKindService
}

func NewCreateMemoService(repo SaveRepo, kids ...SaveKindService) CreateMemoService {
	kinds := make(map[string]SaveKindService)

	for _, k := range kids {
		kinds[k.GetKind()] = k
	}

	return CreateMemoService{repo, kinds}
}

func (s CreateMemoService) Save(ctx context.Context, memo *Memo) (*Memo, error) {
	record, err := s.kinds[memo.Kind].Save(ctx, NewMemo(memo, WithStatus("RECEIVED")))
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveKind(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}
