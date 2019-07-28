package server

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/pelletier/go-toml"
	"socket_server/config"
	"socket_server/utils"
	"time"
)

var errTypeInt64 = errors.New("类型转换错误: 未能将 interface{} 转换为 int64 类型")

var (
	Store = New()
)

type RedisStore struct {
	pool *redis.Pool
}

func New() *RedisStore {
	redisConfig := config.Config.Get("redis").(*toml.Tree)
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(redisConfig.Get("network").(string), redisConfig.Get("address").(string))
			if err != nil {
				return nil, err
			}
			password := redisConfig.Get("password").(string)
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedisStore{pool}
}

func (r *RedisStore) Find(userId string) (sessionId string, exists bool, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	s, err := redis.String(conn.Do("GET", userId))
	if err == redis.ErrNil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	return s, true, nil
}

func (r *RedisStore) Save(userId string, sessionId string) error {
	conn := r.pool.Get()
	defer conn.Close()
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = conn.Send("SET", userId, sessionId)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXEC")
	return err
}

func (r *RedisStore) EnQueue(msgKey string, mgsInfo string) error {
	conn := r.pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", msgKey, mgsInfo)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) CountQueued(msgKey string) (int32, error) {
	conn := r.pool.Get()
	defer conn.Close()
	queueLength, err := conn.Do("LLEN", msgKey)
	if err != nil {
		return 0, err
	}

	count, ok := queueLength.(int32)
	if !ok {
		return 0, errTypeInt64
	}
	return int32(count), nil
}

func (r *RedisStore) GetQueue(msgKey string) ([]string, error) {
	msgInfoStrings := make([]string, 0)
	conn := r.pool.Get()
	defer conn.Close()
	msgCount, err := r.CountQueued(msgKey)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(msgCount); i++ {
		msgInfoString, err := conn.Do("RPOP", msgKey)
		if err != nil {
			return msgInfoStrings, err
		}
		msgInfoStrings = append(msgInfoStrings, utils.B2String(msgInfoString.([]uint8)))
	}
	return msgInfoStrings, nil
}
