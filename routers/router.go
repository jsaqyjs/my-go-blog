package routers

import (
	"github.com/gin-gonic/gin"
	"my-go-blog/service/article_service"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/articles", func(c *gin.Context) {
		articles := article_service.GetArticleList()
		c.JSON(200, gin.H{
			"message": "查询文章成功",
			"data":    articles,
			"code":    200,
		})
	})
	return r
}
