package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/goccy/go-json"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Memo struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

func withReceivedStatus(memo Memo) Memo {
	memo.Status = "RECEIVED"
	return memo
}

func CreateMemo(url, memoType string) (*Memo, error) {
	endpoints := []string{"http://localhost:2379"}
	dialTimeout := 1 * time.Second

	// Create a new etcd client
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}
	defer cli.Close() // Make sure to close the client

	id := getHex()

	m := newMemo(id, memoType, "RECEIVED", url)

	err = put(cli, m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func newMemo(id, memoType, status, url string) *Memo {
	return &Memo{
		ID:     id,
		Type:   memoType,
		Status: status,
		Url:    url,
	}
}

func put(cli *clientv3.Client, memo *Memo) error {
	id := memo.ID
	url := memo.Url

	b, err := json.Marshal(memo)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	sec, _ := time.ParseDuration("3h")

	resp, err := cli.Grant(ctx, int64(sec.Seconds()))
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, id, string(b), clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, url, id, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	return nil
}

func getHex() string {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		fmt.Println("error")
	}
	var buf [8]byte
	hex.Encode(buf[:], b[:])
	return string(buf[:])
}
