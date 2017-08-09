package db

import (
	"fmt"
	"testing"

	"github.com/gfandada/gserver/logger"
)

/***
CREATE TABLE `userinfo` (
	`uid` INT(10) NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(64) NULL DEFAULT NULL,
	`departname` VARCHAR(64) NULL DEFAULT NULL,
	`created` DATE NULL DEFAULT NULL,
	PRIMARY KEY (`uid`)
);

CREATE TABLE `userdetail` (
	`uid` INT(10) NOT NULL DEFAULT '0',
	`intro` TEXT NULL,
	`profile` TEXT NULL,
	PRIMARY KEY (`uid`)
)
***/

func Test_mysql(t *testing.T) {
	logger.Start("../gservices/test.xml")

	//	//插入数据
	//	stmt, _ := GetMysql().Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	//	res, _ := stmt.Exec("gfandada@gmail.com", "test", "2017-08-06")
	//	id, _ := res.LastInsertId()
	//	fmt.Println(id)
	//	res, _ = stmt.Exec("gfandada@gmail.com1", "test1", "2017-08-07")
	//	id, _ = res.LastInsertId()
	//	fmt.Println(id)

	//	//更新数据
	//	stmt, _ = GetMysql().Prepare("update userinfo set username=? where uid=?")
	//	res, _ = stmt.Exec("astaxieupdate", id)
	//	affect, _ := res.RowsAffected()
	//	fmt.Println(affect)

	//	//查询
	//	rows, _ := GetMysql().Query("SELECT * FROM userinfo")
	//	clos, _ := rows.Columns()
	//	fmt.Println("哈哈", clos)
	//	for rows.Next() {
	//		var uid int
	//		var username string
	//		var department string
	//		var created string
	//		rows.Scan(&uid, &username, &department, &created)
	//		//fmt.Println(uid, username, department, created)
	//	}
	//	res, _ = GetMysql().Exec("SELECT * FROM userinfo where uid=3")
	//	f
	//	for rows.Next() {
	//		var uid int
	//		var username string
	//		var department string
	//		var created string
	//		rows.Scan(&uid, &username, &department, &created)
	//		fmt.Println("恩恩", uid, username, department, created)
	//	}

	//	//删除数据
	//	stmt, _ = GetMysql().Prepare("delete from userinfo where uid=?")
	//	res, _ = stmt.Exec(id)
	//	affect, _ = res.RowsAffected()
	//	fmt.Println(affect)

	//	//事务操作
	//	tx, _ := GetMysql().Begin()
	//	stmt, _ = tx.Prepare("delete from userinfo where uid=?")
	//	res, _ = stmt.Exec(1)
	//	tx.Commit()

	// 使用封装的接口查询
	NewMysql(&Mysql{
		User:         "root",
		Password:     "123456",
		Host:         "192.168.78.130:3306",
		Db:           "gs",
		MaxOpenConns: 16,
		MaxIdleConns: 4,
	})
	a, _ := Query("SELECT * FROM userinfo where uid=3")
	for a.Next() {
		var uid int
		var username string
		var department string
		var created string
		a.Scan(&uid, &username, &department, &created)
		fmt.Println("非事务查询", uid, username, department, created)
	}
	a.Close()
	fmt.Println(Update("update userinfo set username=? where uid=?", "fanlin", 3))
	fmt.Println(Delete("delete  from userinfo where uid=?", 3))
	fmt.Println(Insert("insert into userinfo(username,departname) values(?,?)",
		"xxxxxx", "a"))
	// 事务操作
	trans, err := Begin()
	if err != nil {
		t.Error(err)
	}
	rows, err := trans.Query("SELECT * FROM userinfo")
	var uid int
	var username string
	var department string
	var created string
	for rows.Next() {
		rows.Scan(&uid, &username, &department, &created)
		fmt.Println("事务查询", uid, username, department, created)
	}
	rows.Close()
	if err != nil {
		trans.Rollback()
	} else {
		trans.Commit()
	}
}
