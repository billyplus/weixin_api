package repo

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/billyplus/weixin_api"
	"github.com/pkg/errors"
)

var _ weixin_api.IRepository = (*Memory)(nil)

type Memory struct {
	accessToken        string
	accessTokenExpired time.Time
	mut                int32
}

func (memo *Memory) GetAccessToken(_ context.Context) (string, time.Time, error) {
	if memo.accessTokenExpired.Before(time.Now()) {
		return "", memo.accessTokenExpired, nil
	}
	return memo.accessToken, memo.accessTokenExpired, nil
}

func (memo *Memory) UpdateAccessToken(_ context.Context, tok string, expiredTime time.Time) error {
	memo.accessToken = tok
	memo.accessTokenExpired = expiredTime
	return nil
}

func (memo *Memory) Lock() error {
	if atomic.CompareAndSwapInt32(&memo.mut, 0, 1) {
		return nil
	}
	return errors.New("repo is already locked")
}

func (memo *Memory) UnLock() {
	atomic.CompareAndSwapInt32(&memo.mut, 1, 0)
}

// func (memo *Memory) saveToFile() error {

// }
