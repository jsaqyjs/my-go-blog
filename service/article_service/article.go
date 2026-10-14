package article_service

import "my-go-blog/models"

func GetArticleList() ([]models.Article, error) {
	return models.GetArticles()
}

func GetArticleById(id int) (models.Article, error) {
	return models.GetArticleById(id)
}

func CreateArticle(title, content string) (models.Article, error) {
	return models.CreateArticle(title, content)
}

func UpdateArticle(id int, title, content string) (models.Article, error) {
	return models.UpdateArticle(id, title, content)
}

func DeleteArticle(id int) error {
	return models.DeleteArticle(id)
}
