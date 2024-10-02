package usecase

import (
	domain "github.com/audryus/2dpoint.site/internal/domain/memo"
)

type CreateMemoService interface {
	Create(domain.Memo) (domain.Memo, error)
}

type CreateMemo struct {
	createMemo CreateMemoService
}

func NewCreateMemoUC(createMemo CreateMemoService) CreateMemo {
	return CreateMemo{createMemo}
}

func (u CreateMemo) Create(url, memoType string) (domain.Memo, error) {
	return u.createMemo.Create(domain.Memo{Url: url, Type: memoType})
}
