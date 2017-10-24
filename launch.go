package main

import (
	"log"
	"route"

	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	log.Println("Start Service!")
	app.RegisterView(iris.HTML("./public/views", ".html"))
	app.StaticWeb("/assert/javascript", "./public/javascripts")
	app.StaticWeb("/assert/style", "./public/styles")
	app.Post("/user/login", route.UserLogin)
	app.Post("/user/register", route.UserRegister)
	app.Get("/register_confirm", route.RegisterConfirm)
	app.Post("/article/addnew", route.ArticleAddNew)
	app.Post("/article/{id:int}", route.GetArticle)
	app.Post("/article/update", route.AriticleUpdate)
	app.Run(iris.Addr(":8080"))
}
