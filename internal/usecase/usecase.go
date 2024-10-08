package usecase

import (
	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type Deps struct {
	CreateMemoService memo.CreateMemoService
	FetchMemoService  memo.FetchService
}

type UseCases struct {
	CreateMemo CreateMemo
	GetMemo    GetMemo
}

type Record struct {
	Text *memo.Memo   `json:"text"`
	Urls []*memo.Memo `json:"urls"`
}

func NewRecord() *Record {
	return &Record{
		Urls: make([]*memo.Memo, 0),
	}
}

func NewUseCases(deps Deps) UseCases {
	return UseCases{
		CreateMemo: NewCreateMemoUC(deps.CreateMemoService),
		GetMemo:    NewGetMemoUC(deps.FetchMemoService),
	}
}
