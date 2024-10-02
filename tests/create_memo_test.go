package tests

import (
	"testing"

	domain "github.com/audryus/2dpoint.site/internal/domain/memo"
	"github.com/audryus/2dpoint.site/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type MockCreateRepoEtcd struct {
}

func (m MockCreateRepoEtcd) Put(memo domain.Memo) error {
	return nil
}

type MockGetRepoEtcd struct {
}

func (m MockGetRepoEtcd) Get(id string) (domain.Memo, error) {
	return domain.Memo{
		ID:     "1",
		Url:    "https://www.example.com",
		Status: "RECEIVED",
	}, nil
}

func (m MockGetRepoEtcd) GetByUrl(id string) (domain.Memo, error) {
	return domain.Memo{
		ID:     "1",
		Url:    "https://www.example.com",
		Status: "RECEIVED",
	}, nil
}

type MockGetRepoEtcdNotFound struct {
}

func (m MockGetRepoEtcdNotFound) Get(id string) (domain.Memo, error) {
	var memo domain.Memo
	return memo, &domain.NotFoundError{}
}
func (m MockGetRepoEtcdNotFound) GetByUrl(id string) (domain.Memo, error) {
	var memo domain.Memo
	return memo, &domain.NotFoundError{}
}

func TestPutMemoAlreadyExists(t *testing.T) {
	uc := usecase.NewCreateMemoUC(domain.NewCreateMemoService(MockCreateRepoEtcd{}, MockGetRepoEtcd{}))

	memo, err := uc.Create("https://www.google.com", "URL")

	assert.NoError(t, err)
	assert.NotNil(t, memo.ID)
	assert.Equal(t, "1", memo.ID)
	assert.Equal(t, "https://www.example.com", memo.Url)
	assert.Equal(t, "RECEIVED", memo.Status)
}

func TestPutMemo(t *testing.T) {
	uc := usecase.NewCreateMemoUC(domain.NewCreateMemoService(MockCreateRepoEtcd{}, MockGetRepoEtcdNotFound{}))

	memo, err := uc.Create("https://www.google.com", "URL")

	assert.NoError(t, err)
	assert.NotNil(t, memo.ID)
	assert.Equal(t, "https://www.google.com", memo.Url)
	assert.Equal(t, "RECEIVED", memo.Status)
}

func TestGetMemo(t *testing.T) {
	uc := usecase.NewGetMemoUC(domain.NewGetMemoService(MockGetRepoEtcd{}, MockCreateRepoEtcd{}))
	memo, err := uc.Get("1")

	assert.NoError(t, err)
	assert.Equal(t, "1", memo.ID)
	assert.Equal(t, "https://www.example.com", memo.Url)
	assert.Equal(t, "RECEIVED", memo.Status)
}
