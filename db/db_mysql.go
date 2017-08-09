package db

import (
	"database/sql"
	"fmt"

	"github.com/gfandada/gserver/logger"
	_ "github.com/go-sql-driver/mysql"
)

var (
	_mysql *sql.DB
)

type Mysql struct {
	User         string
	Password     string
	Host         string
	Db           string
	MaxOpenConns int
	MaxIdleConns int
}

type CallBack func() error

func NewMysql(cfg *Mysql) {
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

func GetMysql() *sql.DB {
	return _mysql
}

/****************************非事务操作********************************/

// FIXME 使用后请手动rows.Close()
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return GetMysql().Query(query, args...)
}

func execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return GetMysql().Exec(sqlStr, args...)
}

func Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

func Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

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

func Begin() (*MysqlTransaction, error) {
	var trans = &MysqlTransaction{}
	var err error
	if pingErr := GetMysql().Ping(); pingErr == nil {
		trans.SQLTX, err = GetMysql().Begin()
	}
	return trans, err
}

func (t *MysqlTransaction) Rollback() error {
	return t.SQLTX.Rollback()
}

func (t *MysqlTransaction) Commit() error {
	return t.SQLTX.Commit()
}

func (t *MysqlTransaction) execute(sqlStr string, args ...interface{}) (sql.Result, error) {
	return t.SQLTX.Exec(sqlStr, args...)
}

// FIXME 使用后请手动rows.Close()
func (t *MysqlTransaction) Query(queryStr string, args ...interface{}) (*sql.Rows, error) {
	rows, err := t.SQLTX.Query(queryStr, args...)
	return rows, err
}

func (t *MysqlTransaction) Update(updateStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(updateStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}

func (t *MysqlTransaction) Insert(insertStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(insertStr, args...)
	if err != nil {
		return 0, err
	}
	lastid, err := result.LastInsertId()
	return lastid, err

}

func (t *MysqlTransaction) Delete(deleteStr string, args ...interface{}) (int64, error) {
	result, err := t.execute(deleteStr, args...)
	if err != nil {
		return 0, err
	}
	affect, err := result.RowsAffected()
	return affect, err
}
