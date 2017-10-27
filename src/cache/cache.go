package cache

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
)

var db *sql.DB
var red_client *redis.Client

var userinfo string = "login:"

var stmt_insert *sql.Stmt
var stmt_update *sql.Stmt
var stmt_getbyauthor *sql.Stmt
var stmt_getarticle *sql.Stmt
var stmt_userinfo *sql.Stmt
var stmt_simple_article_info *sql.Stmt

var sql_insert string = "insert into article values(null,?,?,?,?,?)"
var sql_getbyauthor string = "select * from article where author=?"
var sql_updatearticle string = "update article set content=? where id=?"
var sql_getarticle string = "select author,title,subtitle,content from article where id=?"
var sql_userinfo string = "select rid,nickname,email,phone,username from userinfo where rid=?"
var sql_simple_article_info = "select author,title,createtime from article order by createtime desc"

func init() {
	log.Println("inti in cache.go")
	red_client = redis.New()
	err := red_client.Connect("127.0.0.1", 6379)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	red_client.Auth()
	db, err = sql.Open("mysql", "dev:dalizi1992@tcp(127.0.0.1:3306)/jack")
	db.SetMaxOpenConns(20)
	stmt_insert, _ = db.Prepare(sql_insert)
	stmt_update, _ = db.Prepare(sql_updatearticle)
	stmt_getarticle, _ = db.Prepare(sql_getarticle)
	stmt_userinfo, _ = db.Prepare(sql_userinfo)
	stmt_simple_article_info, _ = db.Prepare(sql_simple_article_info)
	if err != nil {
		panic(err.Error())
	}
}
func CacheGenSimpleArticleInfo(key string, fld1 string, val string, fld2 string, val2 string, fld3 string, val3 string) error {
	_, err := red_client.HMSet(key, fld1, val, fld2, val2, fld3, val3)
	return err
}
func CacheGetPwdRid(key string, fld1 string, fld2 string) ([]string, error) {
	log.Println(key)
	return red_client.HMGet(key, fld1, fld2)
}
func CacheSetPwdRid(key string, fld1 string, val1 string, fld2 string, val2 string) error {
	_, err := red_client.HMSet(key, fld1, val1, fld2, val2)
	return err

}
func DbAddNewArticle(author string, title string, subtitle string, content string, createtime string) (int, error) {
	res, err := stmt_insert.Exec(author, title, subtitle, content, createtime)
	if err != nil {
		log.Println("mysql 错误")
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
func DbUpdateArticle(articleid int, content string) error {
	_, err := stmt_update.Exec(content, articleid)
	if err != nil {
		log.Println("mysql 错误")
	}
	return err
}
func DbGetArticleContent(articleid int) (string, string, string, string, error) {
	res := stmt_getarticle.QueryRow(articleid)

	var author string
	var title string
	var subtitle string
	var content string
	res.Scan(&author, &title, &subtitle, &content)

	return author, title, subtitle, content, nil
}
func DbGetSimpleArticleInfo() []map[string]string {
	res, _ := stmt_simple_article_info.Query()
	re := make([]map[string]string, 0)
	var author string
	var title string
	var createtime string
	for res.Next() {
		res.Scan(&author, &title, &createtime)
		temp := make(map[string]string)
		temp["author"] = author
		temp["title"] = title
		temp["createtime"] = createtime
		re = append(re, temp)
	}
	return re

}
func DbGetUserinfoByRid(rrid int) (int, string, string, string, string) {
	rows := stmt_userinfo.QueryRow(rrid)
	var rid int
	var nickname string
	var username string
	var email string
	var phone string
	err := rows.Scan(&rid, &username, &email, &phone, &nickname)
	if err != nil {
		return 0, "", "", "", ""
	}
	return rid, nickname, username, email, phone
}
func CacheCheckExistsEmail(key string) (bool, error) {

	return red_client.Exists(key)
}
func CacheGetNextRid(key string) (int64, error) {
	return red_client.Incr(key)
}
