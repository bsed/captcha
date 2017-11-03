package captcha

import (
	"sync"
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

type redisStore struct {
	sync.RWMutex
	client *redis.Client
	// 超时时间
	expiration time.Duration
}

func NewRedisStore(addr string, password string, db int, expiration time.Duration) Store {
	s := new(redisStore)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB

	})
	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println(pong)
	}
	s.expiration = expiration
	s.client = client
	return s
}


func (s *redisStore) Set(id string, digits []byte) {
	s.Lock()
	err :=s.client.Set("captcha_"+id,string(digits),s.expiration).Err()
	if err!=nil {
		panic(err)
	}
	s.Unlock()
}

func (s *redisStore) Get(id string, clear bool)(digits []byte) {
	val, err := s.client.Get("captcha_"+id).Result()
	if err !=nil {
		digits = []byte(val)
	}

	return
}
