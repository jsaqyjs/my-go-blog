package routers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"my-go-blog/dto"
	"my-go-blog/middleware"
	"my-go-blog/models"
	"my-go-blog/service/article_service"
	"my-go-blog/service/user_service"
	"net/http"
	"strconv"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.DELETE("/articles/:id", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
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
		if ok != nil {
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
	r.PUT("/articles/:id", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		idS := c.Param("id")
		id, err := strconv.Atoi(idS)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "文章ID格式错误",
			})
			return
		}
		var request dto.ArticleRequest
		// 2. 接收 JSON 请求体

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
			request.Author,
		)
		if found != nil {
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
	r.GET("/articles/count", func(c *gin.Context) {

		count, err := article_service.CountArticles()

		if handleDBError(c, err) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "统计成功",
			"data": gin.H{
				"total": count,
			},
		})
	})
	r.POST("/articles", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		var request dto.ArticleRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误，请填写标题和内容",
			})
			return
		}
		article, err := article_service.CreateArticle(
			request.Title,
			request.Content,
			request.Author,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "文章创建成功",
			"data":    article,
		})
	})
	r.GET("/articles", func(c *gin.Context) {
		articles, err := article_service.GetArticleList()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"message": "查询文章成功",
			"data":    articles,
			"code":    200,
		})
	})
	r.GET("/articles/search", func(c *gin.Context) {
		title := c.Query("title")
		if title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请输入搜索标题",
			})
			return
		}
		articles, err := article_service.SearchArticles(title)
		if handleDBError(c, err) {
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "搜索成功",
			"data":    articles,
		})
	})
	r.GET("/articles/author/:author", func(c *gin.Context) {
		author := c.Param("author")
		if author == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "作者不能为空",
			})
			return
		}
		articles, err := article_service.GetArticleByAuthor(author)
		if handleDBError(c, err) {
			return
		}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "搜索成功",
			"data":    articles,
		})
	})
	r.GET("/articles/latest", func(c *gin.Context) {

		article, err := article_service.GetLatestArticle()

		if handleDBError(c, err) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": article,
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
		if handleDBError(c, found) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"data":    article,
			"message": "查询文章成功",
		})
	})
	r.GET("/articles/page", func(c *gin.Context) {
		page, err1 := strconv.Atoi(
			c.DefaultQuery("page", "1"))
		pageSize, err2 := strconv.Atoi(
			c.DefaultQuery("pageSize", "10"))
		if err1 != nil || err2 != nil || page < 1 || pageSize < 1 || pageSize > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "分页参数不合法",
			})
			return
		}
		articles, total, err :=
			article_service.GetArticlesByPage(page, pageSize)
		if handleDBError(c, err) {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"data":      articles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"total_pages": (total + int64(pageSize) - 1) /
				int64(pageSize),
		})
	})

	r.POST("/articles/batch", func(c *gin.Context) {
		var request struct {
			Articles []struct {
				Title   string `json:"title" binding:"required"`
				Content string `json:"content" binding:"required"`
				Author  string `json:"author" binding:"required"`
			} `json:"articles" binding:"required,min=1,max=20,dive"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数出错",
			})
			return
		}
		articles := make([]models.Article, 0, len(request.Articles))
		for _, article := range request.Articles {
			articles = append(articles, models.Article{
				Title:   article.Title,
				Content: article.Content,
				Author:  article.Author,
			})
		}
		err := article_service.CreateArticlesBatch(articles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "批量创建失败，事务已回滚",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"code":    201,
			"message": "批量创建成功",
		})
	})
	// 用户登录
	r.POST("/users/login", func(c *gin.Context) {
		var request dto.LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请输入用户名和密码",
			})
			return
		}

		token, err := user_service.Login(request.Username, request.Password)
		if err != nil {

			// 用户名或密码错误
			if errors.Is(err, user_service.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "用户名或密码错误",
				})
				return
			}

			// 数据库或 Token 生成失败
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "登录服务暂时不可用",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": gin.H{"token": token,
				"token_type": "Bearer",
				"expires_in": 7200},
			"message": "登录成功",
		})
	})
	// 用户注册
	r.POST("/users/register", func(c *gin.Context) {
		var request dto.RegisterRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "用户名或密码格式不正确",
			})
			return
		}
		user, err := user_service.Register(
			request.Username,
			request.Password,
		)
		// 3. 判断注册结果
		if err != nil {

			if errors.Is(err, models.ErrUsernameExists) {
				c.JSON(http.StatusConflict, gin.H{
					"code":    409,
					"message": "用户名已存在",
				})
				return
			}
			if errors.Is(err, user_service.ErrInvalidPassword) {
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "密码长度不合法",
				})
				return
			}
			// 其他错误暂时返回服务器错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "注册失败",
			})
			return
		}
		// 4. 返回成功结果，不返回密码
		c.JSON(http.StatusCreated, gin.H{
			"code":    201,
			"message": "注册成功",
			"data": dto.RegisterResponse{
				ID:       user.ID,
				Username: user.Username,
			},
		})
	})
	r.GET(
		"/users/me",
		middleware.JWTAuthMiddleware(),
		func(c *gin.Context) {
			// 1. 获取中间件保存的用户 ID
			userID, ok := middleware.GetCurrentUserID(c)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "无法获取当前用户",
				})
				return
			}
			// 2. 查询用户
			user, err := user_service.GetCurrentUser(userID)

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code":    401,
						"message": "用户不存在",
					})
					return
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "查询用户失败",
				})
				return
			}

			// 3. 返回当前用户的信息
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"data": gin.H{
					"id":       user.ID,
					"username": user.Username,
				},
			})
		})

	return r
}

func handleDBError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return true
	}

	c.JSON(500, gin.H{
		"code":    500,
		"message": "数据库操作失败",
	})

	return true
}
