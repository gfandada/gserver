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
	NewMysql("./mysql.json")
	ret, err := Query("SELECT * FROM userinfo where uid=5")
	if err != nil {
		t.Error(err)
		return
	}
	for k := range ret {
		fmt.Println("第", k, "行")
		for v := range ret[k] {
			fmt.Println(v, ret[k][v])
		}
	}
	fmt.Println(Update("update userinfo set username=? where uid=?", "fanlin", 3))
	fmt.Println(Delete("delete  from userinfo where uid=?", 3))
	fmt.Println(Insert("insert into userinfo(username,departname) values(?,?)",
		"xxxxxx", "a"))
	// 事务操作
	trans, err := Begin()
	if err != nil {
		t.Error(err)
		return
	}
	ret1, err := trans.Query("SELECT * FROM userinfo")
	if err != nil {
		t.Error(err)
		return
	}
	for k := range ret1 {
		fmt.Println("事务查询第", k, "行")
		for v := range ret1[k] {
			fmt.Println(v, ret1[k][v])
		}
	}
	if err != nil {
		trans.Rollback()
	} else {
		trans.Commit()
	}

	// 来试试事务更新
	trans1, err1 := Begin()
	if err1 != nil {
		t.Error(err1)
		return
	}
	ret2, err := trans1.Query("SELECT * FROM userinfo where uid=59")
	if err != nil {
		t.Error(err)
		return
	}
	if len(ret2) != 1 {
		t.Errorf("1111111")
		trans.Rollback()
		return
	}
	fmt.Println(ret2[0]["created"])
	// 条件不满足不更新
	//	if ret2[0]["created"] != "" {
	//		t.Errorf("222222")
	//		trans1.Rollback()
	//		return
	//	}
	// 更新
	id, err := trans1.Update("UPDATE userinfo SET created=?", "2000-12-12")
	if err != nil {
		t.Errorf("333333")
		trans1.Rollback()
		return
	}
	fmt.Println(id)

	// 正常提交
	// trans1.Commit()
	// 测试回滚
	trans1.Rollback()
}
