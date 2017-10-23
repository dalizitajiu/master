package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
)

var db *sql.DB
var red_client *redis.Client
var stmt_insert *sql.Stmt
var stmt_getbyauthor *sql.Stmt

var sql_insert string = "insert into article values(null,?,?,?,?)"
var sql_getbyauthor string = "select * from article where author=?"

func GetTags(raw string) []string {
	relist := make([]string, 0)
	return append(relist, "sfsf")
}
func AddNewArticle(anthor string, title string, subtitle string, content string) error {
	res, err := stmt_insert.Exec(anthor, title, subtitle, content)
	id, err := res.LastInsertId()

	return err
}

func Init() error {
	stmt_insert, _ = db.Prepare(sql_insert)
	stmt_getbyauthor, _ = db.Prepare(sql_getbyauthor)
	return nil
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
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(timestamp)

	res, err := stmt_insert.Exec("lix", "title3", "subtitle1", "kongneirong")

	if err != nil {
		fmt.Println(err.Error())
	}
	//	defer stmt_insert.Close()
	//	defer stmt_getbyauthor.Close()
	id, err := res.LastInsertId()
	fmt.Println(id)
}
