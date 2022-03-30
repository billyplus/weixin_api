package repo

import (
	"context"
	"time"
)

type Memory struct {
	accessToken        string
	accessTokenExpired time.Time
}

func (memo *Memory) GetAccessToken(_ context.Context) (string, error) {
	if memo.accessTokenExpired.Before(time.Now()) {
		return "", nil
	}
	return memo.accessToken, nil
}

func (memo *Memory) UpdateAccessToken(_ context.Context, tok string, expiredTime time.Time) error {
	memo.accessToken = tok
	memo.accessTokenExpired = expiredTime
	return nil
}

// func (memo *Memory) saveToFile() error {

// }
