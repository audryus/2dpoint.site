package memo

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Memo struct {
	ID      string           `json:"id"`
	Kind    string           `json:"kind"`
	Status  string           `json:"status"`
	Content string           `json:"content"`
	Hash    string           `json:"hash_code"`
	Lease   clientv3.LeaseID `json:"lease_id"`
}

func New(content, kind string) *Memo {
	return &Memo{
		Content: content,
		Kind:    kind,
		Hash:    Hash(content),
	}
}

func Hash(content string) string {
	return fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(content)))
}

func NewMemo(b *Memo, options ...func(*Memo)) *Memo {
	svr := &Memo{
		ID:      b.ID,
		Kind:    b.Kind,
		Status:  b.Status,
		Content: b.Content,
		Hash:    b.Hash,
		Lease:   b.Lease,
	}

	for _, o := range options {
		o(svr)
	}
	return svr
}

const Timeout = 1 * time.Second

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "memo not found"
}

func GenerateID() (string, error) {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		return "", err
	}
	var buf [8]byte
	hex.Encode(buf[:], b[:])
	return string(buf[:]), nil
}

func WithID(id string) func(*Memo) {
	return func(s *Memo) {
		s.ID = id
	}
}

func WithLease(lease clientv3.LeaseID) func(*Memo) {
	return func(s *Memo) {
		s.Lease = lease
	}
}

func WithStatus(status string) func(*Memo) {
	return func(s *Memo) {
		s.Status = status
	}
}

func WithHash(content string) func(*Memo) {
	return func(s *Memo) {
		s.Hash = Hash(content)
	}
}
