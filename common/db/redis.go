package db

type RedisConfig struct {
	//Password     string       `yaml:"Password"` // Redis服务的访问账号密码
	//RedisSimple  RedisSimple  `yaml:"Simple"`   // Redis单机模式配置
	//RedisCluster RedisCluster `yaml:"Cluster"`  // Redis集群模式配置
	Addrs        string `yaml:"Addrs"`        // redis address
	DB           int    `yaml:"DB"`           // Redis database instance ID, default is 0
	Password     string `yaml:"Password"`     // Redis database instance ID, default is 0
	ReadTimeout  int    `yaml:"ReadTimeout"`  // socket read timeout, unit: millisecond, default: 3000
	WriteTimeout int    `yaml:"WriteTimeout"` // socket write timeout, unit: millisecond
	DialTimeout  int    `yaml:"DialTimeout"`  // redis connection timeout, default is 5 second

	PoolSize           int `yaml:"PoolSize"`           // max connections in redis connection pool, default is 10 * CPU cores
	MinIdleConns       int `yaml:"MinIdleConns"`       // min idle connections in redis connection pool
	MaxConnAge         int `yaml:"MaxConnAge"`         // max lifetime of redis connection, unit: second, default is 0 (no limit)
	PoolTimeout        int `yaml:"PoolTimeout"`        // max wait time after getting a connection from pool, unit: second
	IdleTimeout        int `yaml:"IdleTimeout"`        // idle time of redis connection, unit: second, default is 5 minutes. -1 means disable idle connection check
	IdleCheckFrequency int `yaml:"IdleCheckFrequency"` // frequency of idle connection check, unit: second, default is 1 minute. -1 means disable idle connection check

	// Runtime Mode, default is false
	IsClusterMode bool `yaml:"IsClusterMode"`
}
