package cache

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" //初始化
	"github.com/gosexy/redis"
)

var db *sql.DB
var redisClient *redis.Client

var userinfo = "login:"

var stmtInsert *sql.Stmt
var stmtUpdate *sql.Stmt
var stmtGetByAuthor *sql.Stmt
var stmtGetArticle *sql.Stmt
var stmtGetUserInfo *sql.Stmt
var stmtGetSimpleArticleInfo *sql.Stmt

var sqlInsert = "insert into article values(null,?,?,?,?)"
var sqlGetByAuthor = "select * from article where author=?"
var sqlUpdateArticle = "update article set content=? where id=?"
var sqlGetArticle = "select author,title,content,createtime from article where id=?"
var sqlGetUserInfo = "select rid,nickname,email,phone,username from userinfo where rid=?"
var sqlGetSimpleArticleInfo = "select id,author,title,createtime from article order by createtime desc limit ?,?"

func init() {
	log.Println("inti in cache.go")
	redisClient = redis.New()
	err := redisClient.Connect("127.0.0.1", 6379)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	redisClient.Auth()
	db, err = sql.Open("mysql", "dev:dalizi1992@tcp(127.0.0.1:3306)/jack")
	db.SetMaxOpenConns(20)
	stmtInsert, _ = db.Prepare(sqlInsert)
	stmtUpdate, _ = db.Prepare(sqlUpdateArticle)
	stmtGetArticle, _ = db.Prepare(sqlGetArticle)
	stmtGetUserInfo, _ = db.Prepare(sqlGetUserInfo)
	stmtGetSimpleArticleInfo, _ = db.Prepare(sqlGetSimpleArticleInfo)

	if err != nil {
		panic(err.Error())
	}
}

//RedSimpleArticleInfo 获取简单的文章信息
func RedSimpleArticleInfo(key string, fld1 string, val string, fld2 string, val2 string, fld3 string, val3 string) error {
	_, err := redisClient.HMSet(key, fld1, val, fld2, val2, fld3, val3)
	return err
}

//RedGetPwdRid redis得到密码和rid
func RedGetPwdRid(key string, fld1 string, fld2 string) ([]string, error) {
	log.Println(key)
	return redisClient.HMGet(key, fld1, fld2)
}

//RedSetPwdRid redis设置密码和rid
func RedSetPwdRid(key string, fld1 string, val1 string, fld2 string, val2 string) error {
	_, err := redisClient.HMSet(key, fld1, val1, fld2, val2)
	return err

}

//DbAddNewArticle mysql增加新纹章
func DbAddNewArticle(author string, title string, content string, createtime string) (int, error) {
	res, err := stmtInsert.Exec(author, title, content, createtime)
	if err != nil {
		log.Println("mysql 错误")
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

//DbUpdateArticle mysql更新文章
func DbUpdateArticle(articleid int, content string) error {
	_, err := stmtUpdate.Exec(content, articleid)
	if err != nil {
		log.Println("mysql 错误")
	}
	return err
}

//DbGetArticleContent mysql获取文章内容
func DbGetArticleContent(articleid int) (string, string, string, string, error) {
	res := stmtGetArticle.QueryRow(articleid)

	var author string
	var title string
	var content string
	var createtime string
	res.Scan(&author, &title, &content, &createtime)
	log.Println(author, title, content, createtime)
	return author, title, content, createtime, nil
}

//DbGetSimpleArticleInfo mysql获取简单的文章信息
func DbGetSimpleArticleInfo(pageno int, pagenum int) []map[string]string {
	res, _ := stmtGetSimpleArticleInfo.Query(pageno, pagenum)
	re := make([]map[string]string, 0)
	var id string
	var author string
	var title string
	var createtime string
	for res.Next() {
		res.Scan(&id, &author, &title, &createtime)
		temp := make(map[string]string)
		temp["author"] = author
		temp["title"] = title
		temp["createtime"] = createtime
		temp["id"] = id
		re = append(re, temp)
	}
	return re

}

//DbGetUserinfoByRid msql通过rid获取用户信息
func DbGetUserinfoByRid(rrid string) (string, string, string, string, string) {
	rows := stmtGetUserInfo.QueryRow(rrid)
	var rid string
	var nickname string
	var username string
	var email string
	var phone string
	err := rows.Scan(&rid, &username, &email, &phone, &nickname)
	if err != nil {
		return "", "", "", "", ""
	}
	log.Println("in dbgetuserinfobyrid")
	return rid, username, email, phone, nickname
}

//RedCheckExistsEmail redis检测该邮箱有没有被使用
func RedCheckExistsEmail(key string) (bool, error) {

	return redisClient.Exists(key)
}

//RedGetNextRid redis获取下一个rid
func RedGetNextRid(key string) (int64, error) {
	return redisClient.Incr(key)
}
