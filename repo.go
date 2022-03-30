package weixinapi

import (
	"context"
	"time"
)

type IRepository interface {
	GetAccessToken(ctx context.Context) (string, error)
	UpdateAccessToken(ctx context.Context, tok string, expiredTime time.Time) error
}
