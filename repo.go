package weixin_api

import (
	"context"
	"time"
)

type IRepository interface {
	GetAccessToken(ctx context.Context) (string, time.Time, error)
	UpdateAccessToken(ctx context.Context, tok string, expiredTime time.Time) error
	Lock() error // 上锁，并返回 true表示上锁成功
	UnLock()
}
