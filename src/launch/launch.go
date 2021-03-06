package main

import (
	"log"
	"route"

	_ "net/http/pprof"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	customLogger := logger.New(logger.Config{
		Status:            true,
		IP:                true,
		Method:            true,
		Path:              true,
		MessageContextKey: "logger_message",
	})

	app.Use(customLogger)
	log.Println("Start Service!")
	app.RegisterView(iris.HTML("../public/views", ".html"))
	// app.RegisterView(iris.HTML("F://myvue/dist/", ".html"))
	app.StaticWeb("/assert/javascript", "../public/javascripts")
	app.StaticWeb("/assert/style", "../public/styles")

	app.Get("/", route.RenderIndex)
	app.Get("/view/article/{id:int}", route.RenderArticle)
	app.Get("/view/login", route.RenderLogin)
	app.Get("/view/articlelist", route.RenderSideBar)
	app.Get("/view/article/new", route.MiddleArticleNew, route.RenderAddNewArticle)
	app.Get("/view/article/update/{id:int}", route.RenderUpdate)
	app.Get("/view/article/getones", route.RenderGetOnes)

	app.Get("/hello", route.Hello)
	app.Post("/user/login", route.UserLogin)
	app.Post("/user/register", route.UserRegister)
	app.Get("/register_confirm", route.RegisterConfirm)
	app.Post("/article/addnew", route.MiddleAuth, route.ArticleAddNew)
	app.Get("/article/{id:int}", route.GetArticle)
	app.Post("/article/update", route.MiddleAuth, route.AriticleUpdate)
	app.Get("/article/abstractlist", route.GetArticleList)
	app.Get("/article/gettoken", route.MiddleAuth, route.GetArticleToken)
	app.Get("/article/getones", route.MiddleAuth, route.GetArticlesByRid)
	app.Get("/article/getbytype", route.GetArticleByType)
	app.Run(iris.Addr(":8080"))
}
