package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mozillazg/request"
)

type Res struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}
type UserInfo struct {
	Rid      int    `json:"rid"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}
type ArticleInfo struct {
	Author   string `json:"author"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

var db *sql.DB
var red_client *redis.Client

var register_key string = "awec12*"

var new_article_key string = "article:"

const mail_body string = `<html><body><a src="%s"></body></html>`
const prefix_confirm string = "http://127.0.0.1:8080/register_confirm"
const auth_key string = "auth123&"
const article_key string = "article&#$12"

var stmt_insert *sql.Stmt
var stmt_update *sql.Stmt
var stmt_getbyauthor *sql.Stmt
var stmt_getarticle *sql.Stmt
var stmt_userinfo *sql.Stmt

var sql_insert string = "insert into article values(null,?,?,?,?)"
var sql_getbyauthor string = "select * from article where author=?"
var sql_updatearticle string = "update article set content=? where id=?"
var sql_getarticle string = "select author,title,subtitle,content from article where id=?"

func GetMd5(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

//cookie参数序列化
func Serialize(name string, value string, opt map[string]string) string {
	result := make([]string, 0)
	result = append(result, name+"="+value)
	if val, ok := opt["maxAge"]; ok {
		result = append(result, "Max-Age="+val)
	}
	if val, ok := opt["domin"]; ok {
		result = append(result, "Domin="+val)
	}
	if val, ok := opt["path"]; ok {
		result = append(result, "Path="+val)
	}
	if val, ok := opt["expires"]; ok {
		result = append(result, "Expires="+val)
	}
	if _, ok := opt["httpOnly"]; ok {
		result = append(result, "HttpOnly")
	}
	if _, ok := opt["secure"]; ok {
		result = append(result, "Secure")
	}
	return strings.Join(result, ";")
}
func NewRes(errno int, errmsg string, data interface{}) Res {
	return Res{errno, errmsg, data}
}
func NewUserInfo(rid int, email string, nickname string, phone string, username string) UserInfo {
	return UserInfo{rid, email, nickname, phone, username}
}

func NewAritcleInfo(author string, title string, subtitle string, content string) ArticleInfo {
	return ArticleInfo{author, title, subtitle, content}
}

//发送邮件
func SendEmail(to string, msg string) error {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Data = map[string]string{"to": to, "msg": msg}
	resp, err := req.Post("http://127.0.0.1:5000/send_mail")
	defer resp.Body.Close()
	return err
}

//字符串的base64加密
func GetBase64(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

//base64字符串的解码
func DecodeBase64(raw string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(raw)
	return string(data), err
}

//生成确认邮件的参数
func GetConfirmUrl(email string, password string) string {
	t := time.Now()
	timestamp := strconv.Itoa(int(t.Unix()))
	ecd_email := GetBase64(email)
	ecd_password := GetBase64(password)

	reslist := make([]string, 0)
	reslist = append(reslist, "t="+timestamp)
	reslist = append(reslist, "encoded_mail="+ecd_email)
	reslist = append(reslist, "encoded_password="+ecd_password)
	log.Println(register_key + "|" + timestamp + "|" + email + "|" + password)
	reslist = append(reslist, "token="+GetMd5(register_key+"|"+timestamp+"|"+email+"|"+password))

	return prefix_confirm + "?" + strings.Join(reslist, "&")
}

func GetTags(raw string) []string {
	relist := make([]string, 0)
	return append(relist, "sfsf")
}

//获取当前的时间戳
func GetNow() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

//ruc token校验
func Auth1(ri string, rtime string, rtoken string) bool {
	return GetMd5(auth_key+ri+rtime) == rtoken
}

//文章更新token校验
func Auth2(ri string, articleid int, aricletoken string) bool {
	return GetMd5(article_key+ri+strconv.Itoa(articleid)) == aricletoken
}

//增加新文章
func AddNewArticle(author string, title string, subtitle string, content string) error {
	res, err := stmt_insert.Exec(author, title, subtitle, content)
	if err != nil {
		log.Println("mysql 错误")
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("mysql 错误")
	}
	var tlen int
	if len(content) < 100 {
		tlen = len(content) - 1
	} else {
		tlen = 100
	}
	log.Println(tlen)
	raw := strings.Join(GetTags(content[0:tlen]), "_")
	_, err = red_client.HMSet(new_article_key+strconv.Itoa(int(id)), "author", author, "tags", raw, "createdtime", GetNow())
	return err
}
func UpdateArticle(articleid int, content string) error {
	_, err := stmt_update.Exec(content, articleid)
	if err != nil {
		log.Println("mysql 错误")
	}
	return err
}
func GetArticleContent(articleid int) (string, string, string, string, error) {
	res := stmt_getarticle.QueryRow(articleid)

	var author string
	var title string
	var subtitle string
	var content string
	res.Scan(&author, &title, &subtitle, &content)

	return author, title, subtitle, content, nil
}

func main() {
	app := iris.New()
	//连接mysql数据库
	db, err := sql.Open("mysql", "dev:dalizi1992@tcp(127.0.0.1:3306)/jack")
	db.SetMaxOpenConns(20)
	defer db.Close()
	stmt_insert, _ = db.Prepare(sql_insert)
	stmt_update, _ = db.Prepare(sql_updatearticle)
	stmt_getarticle, _ = db.Prepare(sql_getarticle)
	if err != nil {
		panic(err.Error())
	}
	red_client := redis.New()
	err = red_client.Connect("127.0.0.1", 6379)
	if err != nil {
		log.Fatal(err.Error())
	}
	//	red_client.Auth()
	defer red_client.Quit()

	// Load all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	app.RegisterView(iris.HTML("./public/views", ".html"))
	app.StaticWeb("/assert/javascript", "./public/javascripts")
	app.StaticWeb("/assert/style", "./public/styles")
	app.Get("/", func(ctx context.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hello.html
		ctx.View("hello.html")
	})

	app.Get("/user/{id:long}", func(ctx context.Context) {
		userID, _ := ctx.Params().GetInt64("id")
		ctx.Writef("User ID: %d", userID)
	})
	stat_userinfo, err := db.Prepare("select rid,username,nickname,email,phone from userinfo where rid=?")
	red_userinfo := "login:"
	if err != nil {
		log.Fatal(err.Error())
	}
	app.Post("/user/login", func(ctx context.Context) {
		mail := ctx.PostValueTrim("email")
		pwd := ctx.PostValueTrim("pwd")
		if mail == "" || pwd == "" {
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		reslist, err := red_client.HMGet(red_userinfo+GetMd5(mail), "pwd", "rid")
		if err != nil {
			ctx.JSON(NewRes(1002, "内部错误", ""))
			return
		}

		rsrid, rpwd := reslist[1], reslist[0]
		if rpwd != pwd {
			ctx.JSON(NewRes(1003, "账号或者密码错误", ""))
			return
		}
		rrid, _ := strconv.Atoi(rsrid)
		rows := stat_userinfo.QueryRow(rrid)

		var rid int
		var nickname string
		var username string
		var email string
		var phone string
		rows.Scan(&rid, &username, &nickname, &email, &phone)

		ctx.Header("Set-Cookie", Serialize("fuck", "shit", map[string]string{"maxAge": "60", "path": "/", "domin": "127.0.0.1"}))
		ctx.JSON(NewRes(0, "", NewUserInfo(rid, email, username, phone, nickname)))
		return
	})
	app.Post("/user/register", func(ctx context.Context) {
		mail := ctx.PostValue("email")
		passwd := ctx.PostValue("pwd")
		if mail == "" || passwd == "" {
			ctx.JSON(NewRes(1001, "邮箱或者密码错误", ""))
			return
		}
		b1, _ := red_client.Exists(red_userinfo + GetMd5(mail))
		if b1 == true {
			ctx.JSON(NewRes(1002, "该邮箱已经被使用", ""))
			return
		}

		go SendEmail(mail, GetConfirmUrl(mail, passwd))
		//		if err != nil {
		//			log.Println("send mail fail")
		//			log.Println(err)
		//		}
		ctx.JSON(NewRes(0, "邮件已经发送", ""))
	})
	app.Get("/register_confirm", func(ctx context.Context) {
		rt := ctx.FormValue("t")
		rmail := ctx.FormValue("encoded_mail")
		rpwd := ctx.FormValue("encoded_password")
		rtoken := ctx.FormValue("token")
		if rt == "" || rmail == "" || rpwd == "" || rtoken == "" {
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		t1, err := strconv.Atoi(rt)
		if err != nil {
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		if (int(time.Now().Unix()) - t1) > 30*60 {
			ctx.JSON(NewRes(1002, "链接已失效", ""))
			return
		}
		mail, err1 := DecodeBase64(rmail)
		if err1 != nil {
			log.Println(err1.Error())
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		pwd, err2 := DecodeBase64(rpwd)
		if err2 != nil {
			log.Println(err2.Error())
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		log.Println(rt, mail, pwd, "sdfsfd")
		if GetMd5(register_key+"|"+rt+"|"+mail+"|"+pwd) != rtoken {
			log.Println(register_key + "|" + rt + "|" + mail + "|" + pwd)
			ctx.JSON(NewRes(1003, "token校验失败", GetMd5(register_key+"|"+rt+"|"+strings.Trim(mail, "\"")+"|"+strings.Trim(pwd, "\""))))
			return
		}
		b1, err := red_client.Exists("login:" + GetMd5(mail))
		if b1 == true {
			ctx.JSON(NewRes(1005, "此链接已失效", ""))
			return
		}
		rid, err1 := red_client.Incr("regist:nextrid")
		if err1 != nil {
			log.Println(err1.Error())
			ctx.JSON(NewRes(1004, "内部错误", ""))
			return
		}
		_, err2 = red_client.HMSet("login:"+GetMd5(mail), "rid", rid, "pwd", pwd)
		if err2 != nil {
			log.Println(err2.Error())
			ctx.JSON(NewRes(1004, "内部错误", ""))
			return
		}
		ctx.JSON(NewRes(0, "注册成功", ""))

	})
	app.Post("/article/addnew", func(ctx context.Context) {
		rtime := ctx.PostValue("r_time")
		rtoken := ctx.PostValue("r_token")
		rid := ctx.GetCookie("I")
		author := ctx.PostValue("r_author")
		title := ctx.PostValue("r_title")
		subtitle := ctx.PostValue("r_subtitle")
		content := ctx.PostValue("r_content")
		if rid == "" || rtime == "" || rtoken == "" || author == "" || title == "" || content == "" {
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		t1, _ := strconv.Atoi(GetNow())
		t2, _ := strconv.Atoi(rtime)
		if (t1 - t2) > 86400 {
			ctx.JSON(NewRes(800, "token过期", ""))
			return
		}
		if !Auth1(rid, rtime, rtoken) {
			ctx.JSON(NewRes(1002, "token错误", ""))
			return
		}
		err := AddNewArticle(author, title, subtitle, content)
		if err != nil {
			ctx.JSON(NewRes(1003, "此标题已经存在", ""))
			return
		}
		ctx.JSON(NewRes(0, "成功提交", ""))
	})
	app.Post("/article/{id:int}", func(ctx context.Context) {
		articleid, err := ctx.Params().GetInt("id")
		//		_ := ctx.FormValue("r_time")
		//		_ := ctx.FormValue("r_token")
		//		_ := ctx.GetCookie("I")
		//		b1 := Auth1(rid, rtime, rtoken)

		if err != nil || articleid <= 0 {
			ctx.JSON(NewRes(1001, "不存在该页面", ""))
			return
		}
		author, title, subtitle, content, err := GetArticleContent(articleid)
		if err != nil {
			ctx.JSON(NewRes(1002, "内部错误", ""))
			return
		}
		ctx.JSON(NewRes(0, "", NewAritcleInfo(author, title, subtitle, content)))
		return
	})
	app.Post("/article/update", func(ctx context.Context) {
		rtime := ctx.PostValue("r_time")
		rtoken := ctx.PostValue("r_token")
		ri := ctx.GetCookie("I")
		rrarticleid := ctx.PostValue("r_articleid")
		rarticletoken := ctx.PostValue("r_articletoken")
		//		author := ctx.PostValue("r_author")
		//		title := ctx.PostValue("r_title")
		//		subtitle := ctx.PostValue("r_subtitle")
		content := ctx.PostValue("r_content")

		if ri == "" || rtime == "" || rtoken == "" || content == "" || rrarticleid == "" {
			ctx.JSON(NewRes(1001, "参数错误", ""))
			return
		}
		rarticleid, _ := strconv.Atoi(rrarticleid)
		if !Auth1(ri, rtime, rtoken) {
			ctx.JSON(NewRes(1002, "token错误", ""))
			return
		}
		if !Auth2(ri, rarticleid, rarticletoken) {
			ctx.JSON(NewRes(1002, "token错误", ""))
			return
		}
		err := UpdateArticle(rarticleid, content)
		if err != nil {
			ctx.JSON(NewRes(1003, "更新失败", ""))
			return
		}
		ctx.JSON(NewRes(0, "更新成功", ""))
		return
	})
	app.Run(iris.Addr(":8080"))
}
