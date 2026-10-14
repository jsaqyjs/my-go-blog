package routers

import (
	"github.com/gin-gonic/gin"
	"my-go-blog/service/article_service"
	"net/http"
	"strconv"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.DELETE("/articles/:id", func(c *gin.Context) {
		idS := c.Param("id")
		id, err := strconv.Atoi(idS)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文章ID格式错误",
			})
			return
		}
		ok := article_service.DeleteArticle(id)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文章不存在",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
		})
	})
	r.PUT("/articles/:id", func(c *gin.Context) {
		idS := c.Param("id")
		id, err := strconv.Atoi(idS)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文章ID格式错误",
			})
			return
		}
		// 2. 接收 JSON 请求体
		var request struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
			})
			return
		}
		article, found := article_service.UpdateArticle(
			id,
			request.Title,
			request.Content,
		)
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
			"message": "文章修改成功",
		})
	})
	r.POST("/articles", func(c *gin.Context) {
		var request struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误，请填写标题和内容",
			})
			return
		}
		article := article_service.CreateArticle(
			request.Title,
			request.Content,
		)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "文章创建成功",
			"data":    article,
		})
	})
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
