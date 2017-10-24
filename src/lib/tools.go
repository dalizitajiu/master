package lib

import (
	"cache"
	_ "cache"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mozillazg/request"
)

var register_key string = "awec12*"
var auth_key string = "auth123&"
var article_key string = "article&#$12"
var new_article_key string = "article:"
var red_userinfo string = "login:"
var nextrid_key string = "regist:nextrid"

const prefix_confirm string = "http://127.0.0.1:8080/register_confirm"

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

//获取md5
func GetMd5(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
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

//发送邮件
func SendEmail(to string, msg string) error {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Data = map[string]string{"to": to, "msg": msg}
	resp, err := req.Post("http://127.0.0.1:5000/send_mail")
	defer resp.Body.Close()
	return err
}

//增加新文章
func AddNewArticle(author string, title string, subtitle string, content string) error {
	id, err := cache.DbAddNewArticle(author, title, subtitle, content)
	if err != nil {
		return err
	}
	var tlen int
	if len(content) < 100 {
		tlen = len(content) - 1
	} else {
		tlen = 100
	}
	log.Println(tlen)
	raw := strings.Join(GetTags(content[0:tlen]), "_")
	return cache.CacheGenSimpleArticleInfo(new_article_key+strconv.Itoa(int(id)), "author", author, "tags", raw, "createdtime", GetNow())
}
func UpdateArticle(articleid int, content string) error {
	return cache.DbUpdateArticle(articleid, content)
}
func GetArticleContent(articleid int) (string, string, string, string, error) {
	return cache.DbGetArticleContent(articleid)
}
func GetPwdRid(mail string) ([]string, error) {
	return cache.CacheGetPwdRid(red_userinfo+GetMd5(mail), "pwd", "rid")
}
func SetPwdRid(mail string, pwd string, rid int) error {
	return cache.CacheSetPwdRid(red_userinfo+GetMd5(mail), "pwd", pwd, "rid", strconv.Itoa(rid))
}
func GetUserinfoByRid(rid int) (int, string, string, string, string) {
	return cache.DbGetUserinfoByRid(rid)
}
func CheckExistsEmail(email string) bool {
	b1, err := cache.CacheCheckExistsEmail(red_userinfo + GetMd5(email))
	if err != nil {
		return true
	}
	return b1
}
func GetNextRid() (int, error) {
	nextrid, err := cache.CacheGetNextRid(nextrid_key)
	return int(nextrid), err
}
