package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/billyplus/weixin_api"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	cmdSet     = "SET"
	cmdGet     = "GET"
	cmdDel     = "DEL"
	cmdPing    = "PING"
	cmdExpire  = "EXPIRE"
	cmdPexpire = "PEXPIRE"

	// order set
	cmdZadd            = "ZADD"
	cmdZrange          = "ZRANGE"
	cmdZrem            = "ZREM"
	cmdZincrBy         = "ZINCRBY"
	cmdZremRangeByRank = "ZREMRANGEBYRANK"

	// hash map
	cmdHGetAll = "HGETALL"
	cmdHSet    = "HSET"
	cmdHGet    = "HGET"
	cmdHMGet   = "HMGET"
	cmdHDel    = "HDEL"

	// sets
	// add value to keys set
	cmdSadd = "SADD"
	// return number of values in set
	cmdScard = "SCARD"
	// check if a value is menber of set
	cmdSisMenber = "SISMEMBER"
	// pop menbers and remove them from set
	cmdSpop = "SPOP"
	// return random menbers but not remove them
	cmdSrandMenber = "SRANDMEMBER"
	// remove menbers from set
	cmdSrem = "SREM"
)

const (
	keyAccessToken = "WX_API_AccessToken_%s"
	keyLocked      = "WX_API_Repo_Locked_%s"
)

var _ weixin_api.IRepository = (*RedisCache)(nil)

type RedisCache struct {
	appId          string
	keyAccessToken string
	keyLocked      string
	pool           *redis.Pool
}

func NewRedisRepo(appId string, host string, port int, password string) *RedisCache {
	addr := fmt.Sprintf("%s:%d", host, port)
	p := &redis.Pool{
		MaxIdle:     32,
		MaxActive:   0,
		IdleTimeout: 5 * time.Minute,
		// Wait:        true,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr, redis.DialPassword(password)) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	rc := &RedisCache{
		appId:          appId,
		keyAccessToken: fmt.Sprintf(keyAccessToken, appId),
		keyLocked:      fmt.Sprintf(keyLocked, appId),
		pool:           p,
		// tokGen: utils.NewUIDGenerator(uint64(time.Now().UnixNano()) << 32),
	}

	return rc
}

func (rc *RedisCache) Close() {
	if rc.pool != nil {
		rc.pool.Close()
	}
}

func (rc *RedisCache) del(ctx context.Context, key string) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := redis.Int(conn.Do(cmdDel, key))
	if err != nil {
		// log.Error("[Get] failed to get record to redis cache", zap.Error(err))
		return errors.Wrap(err, "Del:")
	}
	return nil
}

func (rc *RedisCache) get(ctx context.Context, key string) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	v, err := conn.Do(cmdGet, key)
	if err != nil {
		// log.Error("[Get] failed to get record to redis cache", zap.Error(err))
		return nil, errors.Wrap(err, "Get:")
	}
	return v, nil
}

func (rc *RedisCache) set(ctx context.Context, param ...interface{}) error {
	conn := rc.pool.Get()
	defer conn.Close()

	_, err := conn.Do(cmdSet, param...)
	if err != nil {
		return errors.Wrap(err, "Set:")
	}
	return nil
}

type tokenData struct {
	Tok    string
	Expire time.Time
}

func (rc *RedisCache) GetAccessToken(ctx context.Context) (string, time.Time, error) {
	data, err := redis.Bytes(rc.get(ctx, rc.keyAccessToken))
	if errors.Is(err, redis.ErrNil) {
		return "", time.Now(), nil
	}
	var v tokenData
	if err = json.Unmarshal(data, &v); err != nil {
		return "", time.Now(), err
	}

	return v.Tok, v.Expire, nil
}

func (rc *RedisCache) UpdateAccessToken(ctx context.Context, tok string, expiredTime time.Time) error {
	dur := time.Until(expiredTime).Milliseconds()
	d := &tokenData{
		Tok:    tok,
		Expire: expiredTime.Add(-60 * time.Second),
	}
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	if err = rc.set(ctx, rc.keyAccessToken, data, "PX", int32(dur)); err != nil {
		//log.Debug().Str("key", mut.key).Err(err).Msg("Lock")
		return err
	}
	return nil
}

func (rc *RedisCache) Lock() error {
	ctx := context.Background()
	if err := rc.set(ctx, rc.keyLocked, 1, "NX", "PX", 5000); err != nil {
		//log.Debug().Str("key", mut.key).Err(err).Msg("Lock")
		return errors.WithMessage(err, "failed to lock repo:")
	}
	return nil
}

func (rc *RedisCache) UnLock() {
	ctx := context.Background()
	if err := rc.del(ctx, rc.keyLocked); err != nil {
		log.Error().Str("key", rc.keyLocked).Err(err).Msg("Failed to unlock repo")
		// return errors.WithMessage(err, "failed to unlock repo:")
	}
}
