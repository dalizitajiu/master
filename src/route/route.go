package route

import (
	"lib"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/context"
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

var register_key string = "awec12*"

func NewRes(errno int, errmsg string, data interface{}) Res {
	return Res{errno, errmsg, data}
}
func NewUserInfo(rid int, email string, nickname string, phone string, username string) UserInfo {
	return UserInfo{rid, email, nickname, phone, username}
}

func NewAritcleInfo(author string, title string, subtitle string, content string) ArticleInfo {
	return ArticleInfo{author, title, subtitle, content}
}
func UserLogin(ctx context.Context) {
	mail := ctx.PostValueTrim("email")
	pwd := ctx.PostValueTrim("pwd")
	if mail == "" || pwd == "" {
		ctx.JSON(NewRes(1001, "参数错误", ""))
		return
	}
	reslist, err := lib.GetPwdRid(mail)
	//	reslist, err := cache.()
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
	rid, nickname, username, email, phone := lib.GetUserinfoByRid(rrid)
	ctx.Header("Set-Cookie", lib.Serialize("fuck", "shit", map[string]string{"maxAge": "60", "path": "/", "domin": "127.0.0.1"}))
	ctx.JSON(NewRes(0, "", NewUserInfo(rid, email, username, phone, nickname)))
	return
}
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
	go lib.SendEmail(mail, lib.GetConfirmUrl(mail, passwd))
	ctx.JSON(NewRes(0, "邮件已经发送", ""))
}
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
	if lib.GetMd5(register_key+"|"+rt+"|"+mail+"|"+pwd) != rtoken {
		log.Println(register_key + "|" + rt + "|" + mail + "|" + pwd)
		ctx.JSON(NewRes(1003, "token校验失败", lib.GetMd5(register_key+"|"+rt+"|"+strings.Trim(mail, "\"")+"|"+strings.Trim(pwd, "\""))))
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
func ArticleAddNew(ctx context.Context) {
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
	err := lib.AddNewArticle(author, title, subtitle, content)
	if err != nil {
		ctx.JSON(NewRes(1003, "此标题已经存在", ""))
		return
	}
	ctx.JSON(NewRes(0, "成功提交", ""))
}
func AriticleUpdate(ctx context.Context) {
	rtime := ctx.PostValue("r_time")
	rtoken := ctx.PostValue("r_token")
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

func GetArticle(ctx context.Context) {
	articleid, err := ctx.Params().GetInt("id")
	if err != nil || articleid <= 0 {
		ctx.JSON(NewRes(1001, "不存在该页面", ""))
		return
	}
	author, title, subtitle, content, err := lib.GetArticleContent(articleid)
	if err != nil {
		ctx.JSON(NewRes(1002, "内部错误", ""))
		return
	}
	ctx.JSON(NewRes(0, "", NewAritcleInfo(author, title, subtitle, content)))
	return
}
