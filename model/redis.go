package model

import (
	"github.com/go-redis/redis"
	"github.com/qiusnay/3dorderquery/util"
)

type RedisObj struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
	PoolSize int    `toml:"pool_size"`
}

var rdconf RedisObj

var Redis *redis.Client

func RedisStart() {
	r := RedisObj{}
	client := r.NewClient()
	Redis = client
	// defer client.Close()
}

// 新建client
func (r *RedisObj) NewClient() *redis.Client {
	util.Config().Bind("conf", "redis", &rdconf)
	// 创建client
	client := redis.NewClient(&redis.Options{
		Addr:     rdconf.Addr,
		Password: rdconf.Password,
		DB:       rdconf.Db,
		PoolSize: rdconf.PoolSize,
	})

	// 测试是否有效
	// pong, err := client.Ping().Result()

	return client
}
