package article_service

import "my-go-blog/models"

func GetArticleList() ([]models.Article, error) {
	return models.GetArticles()
}
func GetArticleByAuthor(author string) ([]models.Article, error) {
	return models.GetArticleByAuthor(author)
}
func GetArticleById(id int) (models.Article, error) {
	return models.GetArticleById(id)
}
func GetLatestArticle() (models.Article, error) {
	return models.GetLatestArticle()
}

func SearchArticles(title string) ([]models.Article, error) {
	return models.SearchArticles(title)
}
func CreateArticle(title, content, author string) (models.Article, error) {
	return models.CreateArticle(title, content, author)
}

func UpdateArticle(id int, title, content string) (models.Article, error) {
	return models.UpdateArticle(id, title, content)
}

func DeleteArticle(id int) error {
	return models.DeleteArticle(id)
}
