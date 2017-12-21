package lib

import (
	"cache"
	_ "cache" //初始化数据库
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

var registerKey = "awec12*"
var authKey = "auth123&"
var ariticleKey = "article&#$12"
var newArticleKey = "article:"
var redisUserInfoKey = "login:"
var nextRidKey = "regist:nextrid"

const prefixConfirm = "http://127.0.0.1:8080/register_confirm"

//Serialize cookie参数序列化
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

//GetMd5 获取md5
func GetMd5(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

//GetBase64 字符串的base64加密
func GetBase64(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

//DecodeBase64 base64字符串的解码
func DecodeBase64(raw string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(raw)
	return string(data), err
}

//GetConfirmURL 生成确认邮件的参数
func GetConfirmURL(email string, password string) string {
	t := time.Now()
	timestamp := strconv.Itoa(int(t.Unix()))
	ecdEmail := GetBase64(email)
	ecdPassword := GetBase64(password)

	reslist := make([]string, 0)
	reslist = append(reslist, "t="+timestamp)
	reslist = append(reslist, "encoded_mail="+ecdEmail)
	reslist = append(reslist, "encoded_password="+ecdPassword)
	log.Println(registerKey + "|" + timestamp + "|" + email + "|" + password)
	reslist = append(reslist, "token="+GetMd5(registerKey+"|"+timestamp+"|"+email+"|"+password))

	return prefixConfirm + "?" + strings.Join(reslist, "&")
}

//GetTags 获取标签
func GetTags(raw string) []string {
	relist := make([]string, 0)
	return append(relist, "sfsf")
}

//GenToken 生成ruc token
func GenToken(rid string, ctime string) string {
	return GetMd5(authKey + rid + ctime)
}

//GenArticleToken 生成文章的token
func GenArticleToken(rid string, articleid string) string {
	return GetMd5(ariticleKey + rid + articleid)
}

//GetNow 获取当前的时间戳
func GetNow() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

//GetNowInt 获取当前的时间戳
func GetNowInt() int {
	return int(time.Now().Unix())
}

//GetNowNano 获取当前的时间戳,纳秒级
func GetNowNano() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

//Auth1 ruc token校验
func Auth1(rid string, rtime string, rtoken string) bool {
	return GetMd5(authKey+rid+rtime) == rtoken
}

//Auth2 文章更新token校验
func Auth2(rid string, articleid int, articletoken string) bool {
	return GetMd5(ariticleKey+rid+strconv.Itoa(articleid)) == articletoken
}

//SendEmail 发送邮件
func SendEmail(to string, msg string) error {
	c := new(http.Client)
	req := request.NewRequest(c)
	req.Data = map[string]string{"to": to, "msg": msg}
	resp, err := req.Post("http://127.0.0.1:5000/send_mail")
	defer resp.Body.Close()
	return err
}

//AddNewArticle 增加新文章
func AddNewArticle(author string, title string, content string, atype string) (int, error) {
	id, err := cache.DbAddNewArticle(author, title, content, GetNow(), atype)
	if err != nil {
		log.Println("这里发生了错误")
		return 0, err
	}
	var tlen int
	if len(content) < 100 {
		tlen = len(content) - 1
	} else {
		tlen = 100
	}
	log.Println(tlen)
	raw := strings.Join(GetTags(content[0:tlen]), "_")
	return id, cache.RedSimpleArticleInfo(newArticleKey+strconv.Itoa(int(id)), "author", author, "tags", raw, "createdtime", GetNow())
}

//UpdateArticle 更新文章
func UpdateArticle(articleid int, content string) error {
	return cache.DbUpdateArticle(articleid, content)
}

//GetArticleContent 获取文章内容
func GetArticleContent(articleid int) (string, string, string, string, error) {
	return cache.DbGetArticleContent(articleid)
}

//GetPwdRid 获取密码和rid
func GetPwdRid(mail string) ([]string, error) {
	log.Println("GetPwdRid")
	return cache.RedGetPwdRid(redisUserInfoKey+GetMd5(mail), "pwd", "rid")
}

//SetPwdRid 设置密码和rid
func SetPwdRid(mail string, pwd string, rid int) error {
	return cache.RedSetPwdRid(redisUserInfoKey+GetMd5(mail), "pwd", pwd, "rid", strconv.Itoa(rid))
}

//GetUserinfoByRid 通过rid获取用户信息
func GetUserinfoByRid(rid string) (string, string, string, string, string) {
	return cache.DbGetUserinfoByRid(rid)
}

//CheckExistsEmail 检查邮件是否有人使用
func CheckExistsEmail(email string) bool {
	b1, err := cache.RedCheckExistsEmail(redisUserInfoKey + GetMd5(email))
	if err != nil {
		return true
	}
	return b1
}

//GetNextRid 下一个rid
func GetNextRid() (int, error) {
	nextrid, err := cache.RedGetNextRid(nextRidKey)
	return int(nextrid), err
}

//GetSimpleArticleInfo 获取简单的文章信息
func GetSimpleArticleInfo(pageno int) []map[string]string {
	return cache.DbGetSimpleArticleInfo(pageno*10, 10)
}

//TimeStampToUTC 时间戳转为普通时间
func TimeStampToUTC(timestamp string) string {
	t, _ := strconv.ParseInt(timestamp, 10, 64)
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

//GetArticlesByRid 获取某个rid的文章简要信息
func GetArticlesByRid(rid string) []map[string]string {
	author := cache.DbGetAuthorByRid(rid)
	log.Println("in tools", author)
	return cache.DbGetArticlesByAuthor(author)
}

//GetArticleByType 根据type获取列表
func GetArticleByType(stype string) []map[string]string {
	return cache.DbGetArticleByType(stype)
}
