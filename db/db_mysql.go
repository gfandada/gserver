package db

import (
	"database/sql"
	"fmt"

	Loader "github.com/gfandada/gserver/loader"
	"github.com/gfandada/gserver/logger"
	_ "github.com/go-sql-driver/mysql"
)

var (
	_mysql *sql.DB
)

// 配置
type Mysql struct {
	User         string
	Password     string
	Host         string
	Db           string
	MaxOpenConns int
	MaxIdleConns int
}

type CallBack func() error

func NewMysql(path string) {
	cfg := new(Mysql)
	Loader.LoadJson(path, cfg)
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.User, cfg.Password, cfg.Host, cfg.Db)
	_mysql, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Error("mysql {%s} start error {%v}", dataSourceName, err)
		return
	}
	err = _mysql.Ping()
	if err != nil {
		logger.Error("mysql {%s} ping error {%v}", dataSourceName, err)
		return
	}
	_mysql.SetMaxIdleConns(cfg.MaxIdleConns)
	_mysql.SetMaxOpenConns(cfg.MaxOpenConns)
	return
}

func CloseMysql() {
	if _mysql == nil {
		return
	}
	_mysql.Close()
}

// 获取一个mysql实例
func GetMysql() *sql.DB {
	return _mysql
}

/****************************非事务操作********************************/

func execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return GetMysql().Exec(sqlStr, args...)
}

func Query(queryStr string, args ...interface{}) (map[int]map[string]string, error) {
	query, err := GetMysql().Query(queryStr, args...)
	results := make(map[int]map[string]string)
	if err != nil {
		return results, err
	}
	defer query.Close()
	cols, _ := query.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	i := 0
	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}
	return results, nil
}

// 更新
func Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// 插入
func Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

// 删除
func Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

/****************************事务操作********************************/

type MysqlTransaction struct {
	SQLTX *sql.Tx
}

func (t *MysqlTransaction) execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return t.SQLTX.Exec(sqlStr, args...)
}

// 开启事务
func Begin() (*MysqlTransaction, error) {
	var trans = &MysqlTransaction{}
	var err error
	if pingErr := GetMysql().Ping(); pingErr == nil {
		trans.SQLTX, err = GetMysql().Begin()
	}
	return trans, err
}

// 终止事务
func (t *MysqlTransaction) Rollback() error {
	return t.SQLTX.Rollback()
}

// 提交事务
func (t *MysqlTransaction) Commit() error {
	return t.SQLTX.Commit()
}

// 查询
func (t *MysqlTransaction) Query(queryStr string, args ...interface{}) (map[int]map[string]string, error) {
	query, err := t.SQLTX.Query(queryStr, args...)
	results := make(map[int]map[string]string)
	if err != nil {
		return results, err
	}
	defer query.Close()
	cols, _ := query.Columns()
	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
	i := 0
	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return results, err
		}
		row := make(map[string]string)
		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row
		i++
	}
	return results, nil
}

// 更新
func (t *MysqlTransaction) Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

// 插入
func (t *MysqlTransaction) Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

// 删除
func (t *MysqlTransaction) Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}
