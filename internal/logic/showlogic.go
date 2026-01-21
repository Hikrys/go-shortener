// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"Go_SHORTENER/internal/svc"
	"Go_SHORTENER/internal/types"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	Err404 = errors.New("404")
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// todo: add your logic here and delete this line
	//查看短链接  /lusytc ->重定向到真实的链接 比如 hp.cn/8 -> https://www.baidu.com
	//req.ShortUrl = lusytc
	//1.根据短链接查询原始的长连接
	//布隆过滤器
	exist, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Errorw("l.svcCtx.Filter.Exists() failed", logx.LogField{Key: "err", Value: err.Error()})
	}
	//如果不存在短链接，我们直接返回
	if !exist {
		return nil, Err404
	}
	//1.1 查询数据库之前，增加缓存层
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err404
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//检查过期时间
	// 如果 ExpireAt 有效(Valid=true) 且 当前时间 晚于(After) 过期时间 ，那就已经过期了
	if u.ExpireAt.Valid && time.Now().After(u.ExpireAt.Time) {
		// 既然过期了，我们可以顺手起个协程把它删了，或者标记删除
		go func() {
			l.svcCtx.ShortUrlModel.Delete(context.Background(), u.Id)
		}()
		// 2. 返回 404
		return nil, errors.New("链接已失效")
	}
	// 2. 返回长链
	return &types.ShowResponse{LongUrl: u.Lurl.String}, nil
}
