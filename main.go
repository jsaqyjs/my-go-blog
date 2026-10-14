package main

import (
	"log"
	"my-go-blog/database"
	"my-go-blog/models"
	"my-go-blog/routers"
)

func main() {
	
	// 自动创建 articles 数据表
	if err := database.DB.AutoMigrate(&models.Article{}); err != nil {
		log.Fatal(err)
	}
	r := routers.InitRouter()
	if err := r.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}
