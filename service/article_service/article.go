package article_service

import "my-go-blog/models"

func GetArticleList() []models.Article {
	articleList := models.GetArticles()
	return articleList
}
func GetArticleById(id int) (models.Article, bool) {
	Article, ok := models.GetArticle(id)
	return Article, ok
}
func CreateArticle(title, content string) models.Article {
	newArticle := models.CreateArticle(title, content)
	return newArticle
}

func UpdateArticle(id int, title, content string) (models.Article, bool) {
	article, ok := models.UpdateArticle(id, title, content)
	return article, ok
}
func DeleteArticle(id int) bool {
	return models.DeleteArticle(id)
}
