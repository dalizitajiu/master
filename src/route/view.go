package route

import (
	"lib"
	"log"

	"github.com/kataras/iris/context"
)

func RenderLogin(ctx context.Context) {
	log.Println("RenderArticle")
	ctx.ViewData("message", "lixiomeng"+ctx.Params().Get("id"))
	ctx.View("login.html")
}
func RenderSideBar(ctx context.Context) {
	ctx.View("viewarticle.html")
}
func RenderArticle(ctx context.Context) {
	id, _ := ctx.Params().GetInt("id")
	_, author, title, subtitle, content, createtime := lib.GetArticleContent(id)
	ctx.ViewData("author", author)
	ctx.ViewData("title", title)
	ctx.ViewData("subtitle", subtitle)
	ctx.ViewData("content", content)
	ctx.ViewData("createtime", createtime)
	ctx.View("viewarticle.html")
}
