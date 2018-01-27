package main

import(
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"time"
)
type UserInfo struct{
	Id int `gorm:"primary_key:id;AUTO_INCREMENT"`
	Username string `gorm:"username;type:char(35)"`
	Email string `gorm:"email;type:char(50)"`
	Phone string `gorm:"phone"`
	Nickname string `gorm:"nickname"`
	Rid int `gorm:"rid"`
}

type Article struct{
	Id int `gorm:"primary_key:id;AUTO_INCREMENT"`
	Author string `gorm:"author";type:varchar(45)`
	Title string `gorm:"title"`
	Content string `gorm:content;type:text`
	Createtime int64 `gorm:"createtime"`
	Type string `gorm:"type"` 
}

type Product struct {
	ID    uint `gorm:"primary_key:id"`
	Num   int  `gorm:"AUTO_INCREMENT:number"`
	Code  string
	Price uint
	Tag   []Tag     `gorm:"many2many:tag;"`
	Date  time.Time `gorm:"-"`
}

type Email struct {
	ID         int    `gorm:"primary_key:id"`
	UserID     int    `gorm:"not null;index"`
	Email      string `gorm:"type:varchar(100);unique_index"`
	Subscribed bool
}

type Tag struct {
	Name string
}
func main()  {
	db,err:=gorm.Open("mysql","dev:dalizi1992@tcp(127.0.0.1:3306)/jack?charset=utf8&parseTime=true&loc=Local")
	
	if err!=nil{
		fmt.Println(err)
	}
	
	user:=new(UserInfo)
	db.Last(&user)
	// testuser:=UserInfo{Username:"testuser",Email:"345354@qq.com",Rid:1234567,Phone:"13269101861",Nickname:"Liuyingmei"}
	// db.Create(&testuser)
	// db.NewRecord(testuser)
	fmt.Println(user)

	article:=new(Article)
	db.First(&article)
	fmt.Println(article)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "demo_" + defaultTableName
	}

	db.AutoMigrate(&UserInfo{}, &Email{},&Article{})
	defer db.Close()

}