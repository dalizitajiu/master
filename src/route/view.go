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
	author, title, content, createtime, _ := lib.GetArticleContent(id)
	if author == "" {
		ctx.Redirect("/")
	}
	log.Println(author, title, content, createtime)
	ctx.ViewData("author", author)
	ctx.ViewData("title", title)
	ctx.ViewData("content", content)
	ctx.ViewData("createtime", lib.TimeStampToUTC(createtime))
	ctx.View("article.html")
}

//RenderAddNewArticle 返回增加新文章的模板
func RenderAddNewArticle(ctx context.Context) {
	//todo 访问配置

	ctx.View("editarticle.html")
}

//RenderIndex 默认页
func RenderIndex(ctx context.Context) {
	ctx.View("index.html")
}

//RenderUpdate 更新的页面
func RenderUpdate(ctx context.Context) {
	ctx.View("update.html")
}
