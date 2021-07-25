package dao

import (
	"context"
	"database/sql"
	"fmt"
	"go-/errors"
	"go-/src/homework1/model"
)

/**
	应该抛给上层，个人理解：首先，sql.ErrNoRows是标准库中的底层错误，其次，
	dao层是获取数据对象层，是原子业务中通过sql与数据库交互获取数据对象的一层，
	在它上面还有业务层（biz），服务层（service），如果是http服务器，可能在服务层上还有handler（层）和api（层）
	每一层都有可能由于具体业务的不同，出现新的需要关注的报错信息，那么这一层最底层是必须被wrap的
 */
const (
	USERNAME = "..."
	PASSWORD = "..."
	NETWORK  = "tcp"
	SERVER   = "..."
	PORT     = 3306
	DATABASE = "test_db"
)

var conn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

type DaoImpl struct {
}

func(dao *DaoImpl) QueryById(ctx context.Context, id int64) (*model.Model, error) {
	//建立连接
	db, _ := sql.Open("mysql", conn)
	//建立数据库最大连接数
	db.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	db.SetMaxIdleConns(10)
	//验证连接
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "database conn failed")
	}
	querySql := fmt.Sprintf("select * from test_tab where id = ?")
	rows, err := db.Query(querySql, id)
	//此处的err中包含了ErrNoRow
	if err != nil {
		return nil, errors.Wrap(err, "cannot find rows")
	}
	defer rows.Close()

	var model *model.Model
	for rows.Next() {
		rows.Scan(model)
	}
	return model, nil
}
