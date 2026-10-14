package article_service

import "my-go-blog/models"

func GetArticleList() []models.Article {
	articleList := models.GetArticles()
	return articleList
}
