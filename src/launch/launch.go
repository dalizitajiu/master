package main

import (
	"log"
	"route"

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
	app.StaticWeb("/assert/javascript", "../public/javascripts")
	app.StaticWeb("/assert/style", "../public/styles")

	app.Get("/view/artile/{id:int}", route.RenderArticle)
	app.Get("/view/login", route.RenderLogin)
	app.Get("/test/articlelist", route.RenderSideBar)
	app.Get("/hello", route.Hello)
	app.Post("/user/login", route.UserLogin)
	app.Post("/user/register", route.UserRegister)
	app.Get("/register_confirm", route.RegisterConfirm)
	app.Post("/article/addnew", route.ArticleAddNew)
	app.Post("/article/{id:int}", route.GetArticle)
	app.Post("/article/update", route.AriticleUpdate)
	app.Get("/article/abstractlist", route.GetArticleList)

	app.Run(iris.Addr(":8080"))
}
