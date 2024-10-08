package usecase

import (
	"context"
	"regexp"
	"time"

	"github.com/audryus/2dpoint.site/internal/domain/memo"
)

type CreateMemoService interface {
	Create(content, kind string) (memo.Memo, error)
}

type CreateMemo struct {
	memoService memo.CreateMemoService
}

func NewCreateMemoUC(memoService memo.CreateMemoService) CreateMemo {
	return CreateMemo{
		memoService,
	}
}

func (u CreateMemo) Create(content string) (*Record, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer timeout()

	m, err := u.memoService.Save(ctx, memo.New(content, "text"))
	if err != nil {
		return nil, err
	}

	urls, err := u.getUrls(ctx, content)
	if err != nil {
		return nil, err
	}

	r := new(Record)
	r.Text = m
	r.Urls = append(r.Urls, urls...)

	return r, nil
}

func (u CreateMemo) getUrls(ctx context.Context, text string) ([]*memo.Memo, error) {
	urlPattern := `(https?://|www\.)[^\s]+`
	re := regexp.MustCompile(urlPattern)
	matches := re.FindAllString(text, -1)

	for i, match := range matches {
		matches[i] = regexp.MustCompile(`[^\w/]+$`).ReplaceAllString(match, "")
	}

	urls := make([]*memo.Memo, 0)

	for _, match := range matches {
		m, err := u.memoService.Save(ctx, memo.New(match, "url"))
		if err != nil {
			return nil, err
		}
		urls = append(urls, m)
	}
	return urls, nil
}
