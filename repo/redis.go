package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
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
	keyAccessToken = "KeyAccessToken"
)

type RedisCache struct {
	pool *redis.Pool
}

func NewRedisRepo(host string, port int, password string) *RedisCache {
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
		pool: p,
		// tokGen: utils.NewUIDGenerator(uint64(time.Now().UnixNano()) << 32),
	}

	return rc
}

func (rc *RedisCache) Close() {
	if rc.pool != nil {
		rc.pool.Close()
	}
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

	_, err := redis.String(conn.Do(cmdSet, param...))
	if err != nil {
		return errors.Wrap(err, "Set:")
	}
	return nil
}

func (rc *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	return redis.String(rc.get(ctx, key))
	// if err != nil {
	// 	return "", errors.WithMessage(err, "redis.GetString:")
	// }
	// return v, nil
}

func (rc *RedisCache) GetAccessToken(ctx context.Context) (string, error) {
	tok, err := rc.GetString(ctx, keyAccessToken)
	if errors.Is(err, redis.ErrNil) {
		return "", nil
	}
	return tok, err
	// if err != nil {
	// 	return "", errors.WithMessage(err, "redis.GetString")
	// }
	// return tok, nil
}

func (rc *RedisCache) UpdateAccessToken(ctx context.Context, tok string, expiredTime time.Time) error {
	dur := time.Until(expiredTime).Seconds()
	if err := rc.set(ctx, keyAccessToken, tok, "NX", "PX", int32(dur)); err != nil {
		//log.Debug().Str("key", mut.key).Err(err).Msg("Lock")
		return err
	}
	return nil
}
