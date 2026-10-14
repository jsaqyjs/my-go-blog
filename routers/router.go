package routers

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/articles", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "查询文章成功",
			"data": []gin.H{{
				"id":    1,
				"title": "Go语言学习",
			},
				{
					"id":    2,
					"title": "Gin框架入门",
				}},
		})
	})
	return r
}
