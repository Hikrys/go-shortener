// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"Go_SHORTENER/internal/config"
	"Go_SHORTENER/model"
	"Go_SHORTENER/sequence"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ShortUrlMapModel //对应的是Short_url_map这张表
	Sequence      sequence.Sequence      //sequence
	//Sequence      *sequence.Redis //sequence
	ShortUrlBlackList map[string]struct{}
	//布隆过滤器的初始化
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	//把配置文件中的黑名单配置到map中，便于后续判断
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	//基于go-zero自带的基于Redis的布隆过滤器   //这里是由于new有横线 也就是说 go-zero官方更新的版本中，不推荐这个方式了
	//未来可能弃用，所以我们手动构建一个RedisConf,因为Bloom Filter需要一个redis的客户端
	//store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
	//	r.Type = redis.NodeType
	//})
	redisConf := redis.RedisConf{
		Host: c.CacheRedis[0].Host,
		Type: redis.NodeType,
	}
	store := redis.MustNewRedis(redisConf)
	//声明一个bitSet
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))
	svc := &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(conn, c.CacheRedis),
		Sequence:          sequence.NewMySQL(c.Sequence.DSN),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
	// 布隆过滤器执行预热
	go func() {
		ctx := context.Background()
		offset := uint64(0)
		limit := uint64(1000) // 每次查 1000 条
		for {
			surls, err := svc.ShortUrlModel.FindAllSurls(ctx, offset, limit)
			if err != nil {
				logx.Errorw("BloomFilter Warmup failed", logx.LogField{Key: "err", Value: err.Error()})
				break
			}
			if len(surls) == 0 {
				break
			}

			for _, s := range surls {
				if err := svc.Filter.Add([]byte(s)); err != nil {
					logx.Errorw("BloomFilter Add failed", logx.LogField{Key: "err", Value: err.Error()})
				}
			}

			if len(surls) < int(limit) {
				break
			}
			offset += limit
		}
		fmt.Println("布隆过滤器预热完成")
	}()

	return svc
}
