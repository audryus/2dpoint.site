package usecase

import domain "github.com/audryus/2dpoint.site/internal/domain/memo"

type Deps struct {
	CreateMemoService domain.CreateMemoService
	GetMemoService    domain.GetMemoService
}

type UseCases struct {
	CreateMemo CreateMemo
	GetMemo    GetMemo
}

func NewUseCases(deps Deps) UseCases {
	return UseCases{
		CreateMemo: NewCreateMemoUC(deps.CreateMemoService),
		GetMemo:    NewGetMemoUC(deps.GetMemoService),
	}
}
