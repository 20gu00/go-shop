package redis

import (
	"fmt"
	"github.com/go-redis/redis" // /v8
	"user-web/common/setUp/config"
)

var (
	rdb        *redis.Client
	clusterRDB *redis.ClusterClient
	Nil        = redis.Nil //一种错误,找不到
)

// Init 初始化连接
func InitRedis(cfg *config.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.RedisAddr, cfg.RedisPort),
		Password:     cfg.RedisPassword, // no password set
		DB:           cfg.DB,            // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdle,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func RDBClose() {
	_ = rdb.Close()
}

// sentinel哨兵
//func SentinelRedis() {
//	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
//		MasterName:    "master",
//		SentinelAddrs: []string{"127.0.0.1:6379", "192.168.23.10:6379", "192.168.23.239:6379"},
//	})
//	if _, err := rdb.Ping().Result(); err != nil {
//		return
//	}
//}

// clusterRedis
//func ClusterRedis() {
//	clusterRDB = redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs: []string{":9000", ":9001", "9002"},
//	})
//	if _, err := clusterRDB.Ping().Result(); err != nil {
//		return
//	}
//
//}
