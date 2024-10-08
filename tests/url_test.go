package tests

import (
	"context"
	"testing"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/internal/domain/memo/url"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Url_Save(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 60*time.Minute)
	defer timeout()

	l := logger.New()

	conf, err := config.New(l)
	if err != nil {
		assert.Fail(t, "Error while creating config", err)
		return
	}

	etcdClient, err := etcd.New(conf, l)
	if err != nil {
		assert.Fail(t, "Error while creating etcd client", err)
		return
	}

	urlRepo := url.NewUrlRepo(etcdClient)
	urlCreateS := url.NewCreateUrlService(urlRepo)
	memoRepo := memo.NewMemoRepo(etcdClient)

	s := memo.NewCreateMemoService(memoRepo, urlCreateS)

	content := "https://example.com"

	record, err := s.Save(ctx, memo.New(content, "url"))
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "url", record.Kind)
	assert.Equal(t, content, record.Content)

	record = memo.NewMemo(record, memo.WithStatus("SENT"))

	record, err = s.Save(ctx, record)
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "url", record.Kind)
	assert.Equal(t, content, record.Content)
}

func Test_Url_Get(t *testing.T) {
	ctx, timeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer timeout()

	l := logger.New()
	conf, err := config.New(l)
	if err != nil {
		assert.Fail(t, "Error while creating config", err)
		return
	}

	etcdClient, err := etcd.New(conf, l)
	if err != nil {
		assert.Fail(t, "Error while creating etcd client", err)
		return
	}

	urlRepo := url.NewUrlRepo(etcdClient)
	urlCreateS := url.NewCreateUrlService(urlRepo)
	urlGetS := url.NewGetUrlService(urlRepo)
	memoRepo := memo.NewMemoRepo(etcdClient)

	s := memo.NewCreateMemoService(memoRepo, urlCreateS)
	g := memo.NewFetchService(memoRepo, urlGetS)

	content := "http://test.org"
	record, err := s.Save(ctx, memo.New(content, "url"))
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "url", record.Kind)
	assert.Equal(t, content, record.Content)

	record, err = g.Fetch(ctx, record.ID)
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}

	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "url", record.Kind)
	assert.Equal(t, content, record.Content)
}
