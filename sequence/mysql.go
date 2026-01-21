package sequence

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

//建立一个mysql链接，并且执行插入语句
//REPLACE INTO sequence (stub) VALUES ('a');
//SELECT LAST_INSERT_ID();

const sqlReplaceIntoStub = `REPLACE INTO sequence (stub) VALUES ('a');`

type MySQL struct {
	conn sqlx.SqlConn // go-zero 封装的 sqlx 连接对象
}

func NewMySQL(dsn string) *MySQL {
	return &MySQL{
		conn: sqlx.NewMysql(dsn), // 根据 DSN (数据库地址账号密码) 创建连接
	}
}

// 获取下一个序号
func (m *MySQL) Next() (seq uint64, err error) {
	// 预处理 SQL (Prepare)
	// 这是一个优化习惯，先把 SQL 语句发给 MySQL 让它准备好，防止 SQL 注入（虽然这里没参数）
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw("conn.Prepare failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	defer stmt.Close()
	//执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec() failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	//获取插入主键ID
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId() failed", logx.LogField{Key: "err", Value: err.Error()})
		return 0, err
	}
	return uint64(lid), nil

}
