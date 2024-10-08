package usecase

import (
	"context"
	"time"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type GetMemo struct {
	memoService memo.FetchService
}

func NewGetMemoUC(memoService memo.FetchService) GetMemo {
	return GetMemo{
		memoService,
	}
}

func (u GetMemo) Get(id string) (*memo.Memo, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer timeout()

	return u.memoService.Fetch(ctx, id)
}
