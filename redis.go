package dist

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var (
	ErrNoRedisConf = errors.New("no redis addrs config in zookeeper")
)

// 通用配置
type ZkRedisConf struct {
	Addrs    []string `json:"addrs"`
	Password string   `json:"password"`
	Db       int      `json:"db"`
}

// @data: redis 连接参数 json，来自 zookeeper
// @conf
func ConnectRedis(data []byte) (*redis.Client, error) {
	conf := &ZkRedisConf{}
	err := json.Unmarshal(data, conf)
	if nil != err {
		return nil, err
	}
	// 连接 redis
	//TODO update to cluster version
	rds_addr_len := len(conf.Addrs)
	if rds_addr_len <= 0 {
		return nil, ErrNoRedisConf
	}
	rds_addr := conf.Addrs[0]
	//redis.NewClusterClient(redis.ClusterOptions{})
	conn := redis.NewClient(&redis.Options{
		Addr:         rds_addr,
		Password:     conf.Password,
		DB:           conf.Db,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		PoolTimeout:  60 * time.Second,
	})
	_, err = conn.Ping().Result()
	if nil != err {
		return nil, err
	}
	return conn, nil
}
