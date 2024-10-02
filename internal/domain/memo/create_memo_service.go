package domain

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

type CreateMemoRepoEtcd interface {
	Put(memo Memo) error
}

type CreateMemoService struct {
	createEtcd CreateMemoRepoEtcd
	getEtcd    GetMemoRepoEtcd
}

func NewCreateMemoService(createEtcd CreateMemoRepoEtcd, getEtcd GetMemoRepoEtcd) CreateMemoService {
	return CreateMemoService{
		createEtcd: createEtcd,
		getEtcd:    getEtcd,
	}
}

func (s CreateMemoService) Create(memo Memo) (Memo, error) {
	var m Memo
	oldMemo, err := s.getEtcd.GetByUrl(memo.Url)

	if err != nil {
		_, ok := err.(*NotFoundError)
		if !ok {
			return m, err
		}
	} else {
		return oldMemo, nil
	}

	id, err := getHex()
	if err != nil {
		return m, err
	}

	newMemo := withStatus(withID(memo, id), "RECEIVED")
	return newMemo, s.createEtcd.Put(newMemo)
}

func withStatus(memo Memo, status string) Memo {
	return Memo{
		ID:     memo.ID,
		Type:   memo.Type,
		Status: status,
		Url:    memo.Url,
	}
}

func withID(memo Memo, id string) Memo {
	return Memo{
		ID:   id,
		Type: memo.Type,
		Url:  memo.Url,
	}
}

func getHex() (string, error) {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return "", err
	}
	var buf [8]byte
	hex.Encode(buf[:], b[:])
	return string(buf[:]), nil
}
