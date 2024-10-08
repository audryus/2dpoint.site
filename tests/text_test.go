package tests

import (
	"context"
	"testing"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/internal/domain/memo/text"
	"github.com/audryus/2dpoint.site/pkg/database/cockroach"
	"github.com/audryus/2dpoint.site/pkg/database/etcd"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Text_Save(t *testing.T) {
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
	db, err := cockroach.New(conf, l)
	if err != nil {
		assert.Fail(t, "Error while creating etcd client", err)
		return
	}

	textRepo := text.NewTextRepo(etcdClient, db)
	textCreate := text.NewCreateTextService(textRepo)

	memoRepo := memo.NewMemoRepo(etcdClient)

	s := memo.NewCreateMemoService(memoRepo, textCreate)

	content := "Lorem Ipsum is simply dummy text of the"

	record, err := s.Save(ctx, memo.New(content, "text"))
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "text", record.Kind)
	assert.Equal(t, content, record.Content)

	record = memo.NewMemo(record, memo.WithStatus("SENT"))

	record, err = s.Save(ctx, record)
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "text", record.Kind)
	assert.Equal(t, content, record.Content)
}

func Test_Text_Get(t *testing.T) {
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

	db, err := cockroach.New(conf, l)
	if err != nil {
		assert.Fail(t, "Error while creating etcd client", err)
		return
	}

	textRepo := text.NewTextRepo(etcdClient, db)
	textCreate := text.NewCreateTextService(textRepo)
	textGet := text.NewGetTextService(textRepo)
	memoRepo := memo.NewMemoRepo(etcdClient)

	s := memo.NewCreateMemoService(memoRepo, textCreate)
	g := memo.NewFetchService(memoRepo, textGet)

	content := "It is a long established fact that"
	record, err := s.Save(ctx, memo.New(content, "text"))
	if err != nil {
		assert.Fail(t, "Error while inserting text", err)
		return
	}
	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "text", record.Kind)
	assert.Equal(t, content, record.Content)

	record, err = g.Fetch(ctx, record.ID)
	if err != nil {
		assert.Fail(t, "Error while fetching text", err)
		return
	}

	assert.NotNil(t, record)
	assert.Equal(t, "RECEIVED", record.Status)
	assert.Equal(t, "text", record.Kind)
	assert.Equal(t, content, record.Content)
}
