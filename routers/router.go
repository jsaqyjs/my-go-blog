package routers

import (
	"github.com/gin-gonic/gin"
	"my-go-blog/service/article_service"
	"net/http"
	"strconv"
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

	r.GET("/articles/:id", func(c *gin.Context) {
		idS := c.Param("id")
		id, err := strconv.Atoi(idS)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文章ID格式错误",
			})
			return
		}
		article, found := article_service.GetArticleById(id)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "文章不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    article,
			"message": "查询文章成功",
		})
	})
	return r
}
