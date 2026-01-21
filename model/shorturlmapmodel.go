package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShortUrlMapModel = (*customShortUrlMapModel)(nil)

type (
	ShortUrlMapModel interface {
		shortUrlMapModel
		//预热方法接口
		FindAllSurls(ctx context.Context, offset, limit uint64) ([]string, error)
	}

	customShortUrlMapModel struct {
		*defaultShortUrlMapModel
	}
)

func NewShortUrlMapModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ShortUrlMapModel {
	return &customShortUrlMapModel{
		defaultShortUrlMapModel: newShortUrlMapModel(conn, c, opts...),
	}
}

func (m *customShortUrlMapModel) FindAllSurls(ctx context.Context, offset, limit uint64) ([]string, error) {
	var surls []string
	//这里不走缓存，直接查库。只查 surl 字段，减少 I/O
	query := fmt.Sprintf("select surl from %s where is_del=0 limit ?, ?", m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &surls, query, offset, limit)

	return surls, err
}
