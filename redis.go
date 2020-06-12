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
func ConnectRedisCluster(addrs []string, password string, db int32) (*redis.ClusterClient, error) {
	// 连接 redis
	conn := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     password,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		ReadOnly:     true,
		PoolSize:     10,
		PoolTimeout:  60 * time.Second,
	})
	//conn := redis.NewClient(&redis.Options{
	//	Addr:         rds_addr,
	//	Password:     password,
	//	DB:           int(db),
	//	DialTimeout:  3 * time.Second,
	//	ReadTimeout:  3 * time.Second,
	//	WriteTimeout: 5 * time.Second,
	//	PoolSize:     10,
	//	PoolTimeout:  60 * time.Second,
	//})
	_, err := conn.Ping().Result()
	if nil != err {
		return nil, err
	}
	return conn, nil
}

func ConnectRedis(addr string, password string, db int32) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:         addr,
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
