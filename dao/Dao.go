package dao

import (
	"breplies/data"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DaoHelper struct {
	db *sql.DB
}

func (helper *DaoHelper) New() {
	db, e := sql.Open("mysql", "root:dqndqn@tcp(localhost:3306)/bcomment?charset=utf8")
	if e != nil {
		println(e.Error())
	}
	helper.db = db
}

func (helper *DaoHelper) SaveComment(comment data.Reply) {
	sql := "insert into reply_info(`rpid`,`oid`,`root`,`parent`,`ctime`,`like`,`user_name`,`user_mid`,`content_message`,`content_plat`) values (?,?,?,?,?,?,?,?,?,?)"
	stmt, _ := helper.db.Prepare(sql)
	defer func() {
		if ok := recover(); ok != nil {
			fmt.Println("recover：" + comment.Content.Message)
		}
	}()
	_, e := stmt.Exec(
		comment.Rpid,
		comment.Oid,
		comment.Root,
		comment.Parent,
		comment.Ctime,
		comment.Like,
		comment.Member.Uname,
		comment.Member.Mid,
		comment.Content.Message,
		comment.Content.Plat,
	)
	if e != nil {
		println(e.Error())
	}
}
