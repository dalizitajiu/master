package route

import (
	"log"

	"github.com/kataras/iris/context"
)

func RenderArticle(ctx context.Context) {
	log.Println("RenderArticle")
	ctx.ViewData("message", "lixiomeng"+ctx.Params().Get("id"))
	ctx.View("login.html")
}
