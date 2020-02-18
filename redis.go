package dist

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var (
	ErrInvalidParams = errors.New("invalid parameters")
)

// @data: redis 连接参数 json，来自 zookeeper
// @conf
func ConnectRedis(addrs []string, password string, db int32) (*redis.Client, error) {
	// 连接 redis
	//TODO update to cluster version
	rds_addr_len := len(addrs)
	if rds_addr_len <= 0 {
		return nil, ErrInvalidParams
	}
	rds_addr := addrs[0]
	//redis.NewClusterClient(redis.ClusterOptions{})
	conn := redis.NewClient(&redis.Options{
		Addr:         rds_addr,
		Password:     password,
		DB:           int(db),
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		PoolTimeout:  60 * time.Second,
	})
	_, err := conn.Ping().Result()
	if nil != err {
		return nil, err
	}
	return conn, nil
}
