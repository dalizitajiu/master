package route

import (
	"lib"
	"log"

	"github.com/kataras/iris/context"
)

//RenderLogin 登录页面
func RenderLogin(ctx context.Context) {
	log.Println("RenderArticle")
	ctx.ViewData("message", "lixiomeng"+ctx.Params().Get("id"))
	ctx.View("login.html")
}

//RenderSideBar 侧边栏
func RenderSideBar(ctx context.Context) {
	ctx.View("viewarticle.html")
}

//RenderArticle 文章浏览
func RenderArticle(ctx context.Context) {
	id, _ := ctx.Params().GetInt("id")
	log.Println("当前要浏览的文章id是", id)
	author, title, subtitle, content, createtime, _ := lib.GetArticleContent(id)
	log.Println(author, title, subtitle, content, createtime)
	ctx.ViewData("author", author)
	ctx.ViewData("title", title)
	ctx.ViewData("subtitle", subtitle)
	ctx.ViewData("content", content)
	ctx.ViewData("createtime", createtime)
	ctx.View("article.html")
}
