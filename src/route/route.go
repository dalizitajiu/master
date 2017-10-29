package route

import (
	"lib"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/context"
)

//Res 结果
type Res struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

//UserInfo 用户信息
type UserInfo struct {
	Rid      string `json:"rid"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}

//ArticleInfo 文章信息
type ArticleInfo struct {
	Author     string `json:"author"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"createtime"`
}

//AbstractInfo 摘要信息
type AbstractInfo struct {
	Author     string `json:"author"`
	Title      string `json:"title"`
	CreateTime string `json:"createtime"`
}

var registerKey = "awec12*"

//NewRes 新回复
func NewRes(errno int, errmsg string, data interface{}) Res {
	return Res{errno, errmsg, data}
}

//NewUserInfo 新的用户信息
func NewUserInfo(rid string, email string, nickname string, phone string, username string) UserInfo {
	return UserInfo{rid, email, nickname, phone, username}
}

//NewAritcleInfo 新的文章信息
func NewAritcleInfo(author string, title string, content string, createtime string) ArticleInfo {
	return ArticleInfo{author, title, content, createtime}
}

//Hello 测试用
func Hello(ctx context.Context) {
	ctx.Writef("dsfsd")
}

//UserLogin 登录
func UserLogin(ctx context.Context) {
	log.Println("userinfo")

	mail := ctx.PostValueTrim("email")
	pwd := ctx.PostValueTrim("pwd")
	log.Println(mail, pwd)
	if mail == "" || pwd == "" {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	reslist, err := lib.GetPwdRid(mail)
	log.Println(reslist)
	//	reslist, err := cache.()
	if err != nil {
		ctx.JSON(NewRes(1002, "内部错误", ""))
		return
	}

	rsrid, rpwd := reslist[1], reslist[0]
	log.Println(rsrid, rpwd)
	if rpwd != pwd {
		ctx.JSON(NewRes(1003, "账号或者密码错误", ""))
		return
	}
	rrid := rsrid
	rid, nickname, email, phone, username := lib.GetUserinfoByRid(rrid)
	currenttime := lib.GetNow()
	ctx.Header("Set-Cookie", lib.Serialize("I", rid, map[string]string{"maxAge": "3600*1", "path": "/", "domin": "127.0.0.1"}))
	ctx.Header("Set-Cookie", lib.Serialize("r_time", currenttime, map[string]string{"maxAge": "3600", "path": "/", "domin": "127.0.0.1"}))
	ctx.Header("Set-Cookie", lib.Serialize("r_token", lib.GenToken(rrid, currenttime), map[string]string{"maxAge": "3600", "path": "/", "domin": "127.0.0.1"}))
	ctx.JSON(NewRes(0, "", NewUserInfo(rid, email, username, phone, nickname)))
	return
}

//UserRegister 用户注册
func UserRegister(ctx context.Context) {
	mail := ctx.PostValue("email")
	passwd := ctx.PostValue("pwd")
	if mail == "" || passwd == "" {
		ctx.JSON(NewRes(1001, "邮箱或者密码错误", ""))
		return
	}
	b1 := lib.CheckExistsEmail(mail)
	if b1 == true {
		ctx.JSON(NewRes(1002, "该邮箱已经被使用", ""))
		return
	}
	go lib.SendEmail(mail, lib.GetConfirmURL(mail, passwd))
	ctx.JSON(NewRes(0, "邮件已经发送", ""))
}

//RegisterConfirm 注册邮件confirm
func RegisterConfirm(ctx context.Context) {
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
	mail, err1 := lib.DecodeBase64(rmail)
	if err1 != nil {
		log.Println(err1.Error())
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	pwd, err2 := lib.DecodeBase64(rpwd)
	if err2 != nil {
		log.Println(err2.Error())
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	log.Println(rt, mail, pwd, "sdfsfd")
	if lib.GetMd5(registerKey+"|"+rt+"|"+mail+"|"+pwd) != rtoken {
		log.Println(registerKey + "|" + rt + "|" + mail + "|" + pwd)
		ctx.JSON(NewRes(1003, "token校验失败", lib.GetMd5(registerKey+"|"+rt+"|"+strings.Trim(mail, "\"")+"|"+strings.Trim(pwd, "\""))))
		return
	}
	b1 := lib.CheckExistsEmail(mail)
	if b1 == true {
		ctx.JSON(NewRes(1005, "此链接已失效", ""))
		return
	}
	rid, err1 := lib.GetNextRid()

	if err1 != nil {
		log.Println(err1.Error())
		ctx.JSON(NewRes(1004, "内部错误", ""))
		return
	}
	err2 = lib.SetPwdRid(mail, pwd, rid)
	if err2 != nil {
		log.Println(err2.Error())
		ctx.JSON(NewRes(1004, "内部错误", ""))
		return
	}
	ctx.JSON(NewRes(0, "注册成功", ""))
}

//ArticleAddNew 增加新文章
func ArticleAddNew(ctx context.Context) {
	rid := ctx.GetCookie("I")
	_, _, _, _, nickname := lib.GetUserinfoByRid(rid)
	log.Println("昵称是", nickname)
	author := nickname
	title := ctx.PostValue("r_title")
	log.Println(ctx)
	content := ctx.PostValue("r_content")
	log.Println("昵称是", author, title, content)
	if author == "" || title == "" || content == "" {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}

	id, err := lib.AddNewArticle(author, title, content)
	if err != nil {
		ctx.JSON(NewRes(1003, "此标题已经存在", ""))
		return
	}
	ctx.JSON(NewRes(0, "成功提交", id))
}

//MiddleAuth ruc中间件
func MiddleAuth(ctx context.Context) {
	rid := ctx.GetCookie("I")
	rtime := ctx.GetCookie("r_time")
	rtoken := ctx.GetCookie("r_token")

	if rid == "" || rtime == "" || rtoken == "" {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	t1, _ := strconv.Atoi(lib.GetNow())
	t2, _ := strconv.Atoi(rtime)
	if (t1 - t2) > 86400 {
		ctx.JSON(NewRes(800, "token过期", ""))
		return
	}
	if !lib.Auth1(rid, rtime, rtoken) {
		ctx.JSON(NewRes(1002, "token错误", ""))
		return
	}
	log.Println("RUC校验成功")
	ctx.Next()
}

//AriticleUpdate 文章更新
func AriticleUpdate(ctx context.Context) {
	rtime := ctx.GetCookie("r_time")
	rtoken := ctx.GetCookie("r_token")
	ri := ctx.GetCookie("I")
	rrarticleid := ctx.PostValue("r_articleid")
	rarticletoken := ctx.PostValue("r_articletoken")
	content := ctx.PostValue("r_content")

	if ri == "" || rtime == "" || rtoken == "" || content == "" || rrarticleid == "" {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	rarticleid, _ := strconv.Atoi(rrarticleid)
	if !lib.Auth1(ri, rtime, rtoken) {
		ctx.JSON(NewRes(1002, "token错误", ""))
		return
	}
	if !lib.Auth2(ri, rarticleid, rarticletoken) {
		ctx.JSON(NewRes(1002, "token错误", ""))
		return
	}
	err := lib.UpdateArticle(rarticleid, content)
	if err != nil {
		ctx.JSON(NewRes(1003, "更新失败", ""))
		return
	}
	ctx.JSON(NewRes(0, "更新成功", ""))
	return
}

//GetArticle 获取文章
func GetArticle(ctx context.Context) {
	articleid, err := ctx.Params().GetInt("id")
	if err != nil || articleid <= 0 {
		ctx.JSON(NewRes(1001, "不存在该页面", ""))
		return
	}
	author, title, content, createtime, err := lib.GetArticleContent(articleid)

	if err != nil {
		ctx.JSON(NewRes(1002, "内部错误", ""))
		return
	}
	ctx.JSON(NewRes(0, "", NewAritcleInfo(author, title, content, lib.TimeStampToUTC(createtime))))
	return
}

//GetArticleList 获取简单的文章列表
func GetArticleList(ctx context.Context) {
	rpageno := ctx.FormValue("pageno")
	pageno, err := strconv.Atoi(rpageno)
	if err != nil || pageno < 0 {
		ctx.JSON(NewRes(0, "", ""))
		return
	}
	re := lib.GetSimpleArticleInfo(pageno)
	ctx.JSON(NewRes(0, "", re))
	return
}

//GetArticleToken 生成文章的token
func GetArticleToken(ctx context.Context) {
	rid := ctx.GetCookie("I")
	articleid := ctx.FormValue("articleid")
	_, err := strconv.Atoi(articleid)
	if err != nil {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	ctx.JSON(NewRes(0, "", lib.GenArticleToken(rid, articleid)))
	return
}
