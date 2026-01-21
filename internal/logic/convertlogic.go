// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"Go_SHORTENER/model"
	"Go_SHORTENER/pkg/base62"
	"Go_SHORTENER/pkg/connect"
	"Go_SHORTENER/pkg/md5"
	"Go_SHORTENER/pkg/urltool"
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"

	"Go_SHORTENER/internal/svc"
	"Go_SHORTENER/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链 ： 输入一个长链接  转为 一个短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	//1.校验输入的数据
	//1.1 数据不能为空
	//if len(req.LongUrl) == 0 {}  这里用的是len 而不是 ==nil 是因为其nil是一个空也就是什么都没有，而req是一个
	//结构体指针的类型，就算里面什么都也没有，也不会是nil 而是一个指向空结构体的指针,需要校验的是字段内容，不是对象本身
	//不过这里你也可以直接用validator包来做参数校验 我已经在converthandler完成了数据校验
	//1.2 输入的长链接数据是有效的，可以请求的通的
	//URL的有效性检测
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("无效的链接")
	}
	//1.3 判断转链是否以及存在（查询数据库是否以及存在该长链接）
	//1.3.1 给长链接生成MD5的值，
	md5Value := md5.Sum([]byte(req.LongUrl)) //这里我们使用的是pkg中的我们自己封装的md5 不是标准库中的
	//1.3.2 拿MD5去数据库查询 这里我们不拿长链接直接去查，是为了防止数据库压力过大
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return &types.ConvertResponse{
				ShortUrl: l.svcCtx.Config.ShortDoamin + "/" + u.Surl.String,
			}, nil
		}
		logx.Errorw("ShortUrlModel.FindOneByMd5 failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//1.4 输入的不能已经是一个短链接
	//1.4.1 输入的是一个完整的url hp/1qwe2
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetPath failed", logx.LogField{Key: "lurl", Value: req.LongUrl})
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("该链接已经为短链接")
		}
		logx.Errorw("ShortUrlModel.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	var short string
	for {
		//2.转链接，取号（我们用的取号器）
		//2.1 每次来一个转链请求，我们就使用REPLACE INTO语句往sequence插入一条数据
		//并且取出主键id为号码
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("l.svcCtx.Sequence.Next() failed", logx.LogField{Key: "err", Value: err.Error()})
			return nil, err
		}
		//3.号码转短链接
		//3.1 安全性
		short = base62.Int2String(seq)
		//3.2 短域名避免某些特殊词语
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break //如果生成了黑名单词，就直接退出
		}
	}
	// 【新增】计算过期时间
	var expireAt sql.NullTime
	if req.Seconds > 0 {
		expireAt = sql.NullTime{
			Time:  time.Now().Add(time.Duration(req.Seconds) * time.Second),
			Valid: true,
		}
	} else {
		expireAt = sql.NullTime{Valid: false} // NULL 代表不过期
	}
	//4.保存长短链接到数据库持久化数据
	_, err = l.svcCtx.ShortUrlModel.Insert(l.ctx, &model.ShortUrlMap{
		Lurl:     sql.NullString{String: req.LongUrl, Valid: true},
		Md5:      sql.NullString{String: md5Value, Valid: true},
		Surl:     sql.NullString{String: short, Valid: true},
		ExpireAt: expireAt, //过期时间
	})
	if err != nil {
		// 如果报错是因为 MD5 重复（MySQL错误码 1062），说明刚刚被别人插入了
		// 此时我们应该重新查一遍直接返回
		// 这是一个兜底逻辑
		if isDuplicateEntryError(err) { // 需要自己封装一个判断是否为Duplicate Entry的函数
			u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Value, Valid: true})
			if err == nil {
				return &types.ConvertResponse{
					ShortUrl: l.svcCtx.Config.ShortDoamin + "/" + u.Surl.String,
				}, nil
			}
		}

		logx.Errorw("ShortUrlModel.Insert failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	//4.2 把已经生成的短链接，加载到布隆过滤器
	if err := l.svcCtx.Filter.Add([]byte(short)); err != nil {
		logx.Errorw("Filter.Add failed", logx.LogField{Key: "err", Value: err.Error()})

	}
	//5.返回响应
	//5.1 返回的是短域名 + 短链接 hp/1132t
	shortUrl := l.svcCtx.Config.ShortDoamin + "/" + short
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}

// isDuplicateEntryError 判断是否为数据库唯一键冲突错误
func isDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	// MySQL 的 Duplicate entry 错误通常包含这个字符串
	// 严谨的做法是断言 *mysql.MySQLError 并检查 Number == 1062，但字符串检查在简单场景足够通用
	return strings.Contains(err.Error(), "Duplicate entry")
}
