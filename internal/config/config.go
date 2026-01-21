// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"os"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	BaseString string //base62指定字符串

	ShortUrlBlackList []string

	ShortDoamin string

	CacheRedis cache.CacheConf //redis缓存

	Sequence struct {
		DSN string
	}
}

// 增加一个辅助函数来加载环境变量
func (c *Config) LoadEnv() {
	// 1. 覆盖 MySQL DSN
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		// 假设密码和库名不变，只改 Host (这里为了演示简单，更严谨的做法是把整个 DSN 拼装)
		// 但 DSN 格式复杂 (root:pass@tcp(host:port)/db)，解析替换太麻烦。
		// 所以我们约定：如果环境变量里有 MYSQL_DSN，直接替换整个 DSN
		if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
			c.ShortUrlDB.DSN = dsn
			c.Sequence.DSN = dsn // 假设两个用同一个库
		}
	}

	// 2. 覆盖 Redis Host
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		// CacheRedis 是个列表，通常只配一个
		if len(c.CacheRedis) > 0 {
			c.CacheRedis[0].Host = redisHost
		}
	}
}
