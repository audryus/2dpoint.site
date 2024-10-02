package usecase

import (
	domain "github.com/audryus/2dpoint.site/internal/domain/memo"
)

type GetMemo struct {
	getMemo domain.GetMemoService
}

func NewGetMemoUC(getMemo domain.GetMemoService) GetMemo {
	return GetMemo{getMemo}
}

func (u GetMemo) Get(id string) (domain.Memo, error) {
	return u.getMemo.Get(id)
}
