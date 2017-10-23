package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
)

var db *sql.DB
var red_client *redis.Client
var stmt_insert *sql.Stmt
var stmt_getbyauthor *sql.Stmt
var new_article_key string = "article:"
var sql_insert string = "insert into article values(null,?,?,?,?)"
var sql_getbyauthor string = "select * from article where author=?"

func GetTags(raw string) []string {
	relist := make([]string, 0)
	return append(relist, "sfsf")
}
func GetNow() string {
	return strconv.Itoa(int(time.Now().Unix()))
}
func AddNewArticle(author string, title string, subtitle string, content string) error {
	res, err := stmt_insert.Exec(author, title, subtitle, content)
	if err != nil {
		log.Println("mysql 错误")
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("mysql 错误")
	}
	raw := strings.Join(GetTags(content[0:100]), "_")
	_, err = red_client.HMSet(new_article_key+strconv.Itoa(int(id)), "author", author, "tags", raw, "createdtime", GetNow())
	return err
}

func main() {
	red_client = redis.New()
	red_client.Connect("127.0.0.1", 6379)
	//	red_client.Auth()
	defer red_client.Quit()

	db, err := sql.Open("mysql", "dev:dalizi1992@tcp(127.0.0.1:3306)/jack")

	if err != nil {
		fmt.Println(err.Error())
	}
	db.SetMaxOpenConns(20)
	stmt_insert, _ = db.Prepare(sql_insert)
	stmt_getbyauthor, _ = db.Prepare(sql_getbyauthor)

	res, err := stmt_insert.Exec("lix", "title5", "subtitle1", "kongneirong")

	if err != nil {
		fmt.Println(err.Error())
	}
	log.Println(GetNow())
	id, err := res.LastInsertId()
	fmt.Println(id)
}
