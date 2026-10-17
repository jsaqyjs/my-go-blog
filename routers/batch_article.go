package routers

import (
	"github.com/gin-gonic/gin"
	"my-go-blog/middleware"
	"my-go-blog/service/article_service"
	"net/http"
)

func BatchProcessArticleHandler(c *gin.Context) {
	var req struct {
		ArticleIDs []uint `json:"article_ids" binding:"required,min=1,max=100,dive,gt=0"`
	}
	// 解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文章 ID 参数不合法，最多提交100个",
		})
		return
	}
	// 获取当前登录用户 ID
	userID, ok := middleware.GetCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未登录",
		})
	}
	// 调用 Service
	results := article_service.BatchProcessArticles(
		c.Request.Context(),
		req.ArticleIDs,
		userID,
	)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "文章处理成功",
		"data":    results,
	})
}
