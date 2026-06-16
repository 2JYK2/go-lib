package db

import (
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedisDao(redisConfig *RedisConfig) redis.UniversalClient {
	opts := &redis.UniversalOptions{}

	addr := strings.Split(redisConfig.Addrs, ",")
	if len(addr) != 0 {
		opts.Addrs = addr
	}

	// Database number, default is 0
	if 0 <= redisConfig.DB {
		opts.DB = redisConfig.DB
	}

	//  read timeout, unit: millisecond
	if redisConfig.Password != "" {
		opts.Password = redisConfig.Password
	}

	//  read timeout, unit: millisecond
	if 0 < redisConfig.ReadTimeout {
		opts.ReadTimeout = time.Duration(redisConfig.ReadTimeout) * time.Millisecond
	}
	// write timeout, unit: millisecond
	if 0 < redisConfig.WriteTimeout {
		opts.WriteTimeout = time.Duration(redisConfig.WriteTimeout) * time.Millisecond
	}
	// connection timeout, unit: second
	if 0 < redisConfig.DialTimeout {
		opts.DialTimeout = time.Duration(redisConfig.DialTimeout) * time.Second
	}
	// connection pool size
	if 0 < redisConfig.PoolSize {
		opts.PoolSize = redisConfig.PoolSize
	}
	// minimum idle connections
	if 0 < redisConfig.MinIdleConns {
		opts.MinIdleConns = redisConfig.MinIdleConns
	}
	//  connection maximum lifetime
	if 0 < redisConfig.MaxConnAge {
		opts.ConnMaxLifetime = time.Duration(redisConfig.MaxConnAge) * time.Second
	}
	// maximum wait time for idle connections
	if 0 < redisConfig.PoolTimeout {
		opts.PoolTimeout = time.Duration(redisConfig.PoolTimeout) * time.Second
	}
	// idle connection timeout, unit: minute, -1 means closing the configuration
	if -1 == redisConfig.IdleTimeout {
		opts.ConnMaxIdleTime = time.Duration(-1)
	} else if 0 < redisConfig.IdleTimeout {
		opts.ConnMaxIdleTime = time.Duration(redisConfig.IdleTimeout) * time.Minute
	}
	opts.IsClusterMode = redisConfig.IsClusterMode
	return redis.NewUniversalClient(opts)
}
