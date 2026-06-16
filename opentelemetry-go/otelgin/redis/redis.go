package redis

import (
	"github.com/2JYK2/go-lib/opentelemetry-go/otelgin/redis/redisotel"
	"github.com/redis/go-redis/v9"
)

func RedisTrace(rdb redis.UniversalClient, opts ...redisotel.TracingOption) error {
	return redisotel.InstrumentTracing(rdb, opts...)
}
