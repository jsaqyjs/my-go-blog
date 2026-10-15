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
func CreateArticle(title, content, author string, userId uint) (models.Article, error) {
	return models.CreateArticle(title, content, author, userId)
}

func CountArticles() (int64, error) {
	return models.CountArticles()
}
func UpdateArticle(id int, title, content, author string, userId uint) (models.Article, error) {
	return models.UpdateArticle(id, title, content, author, userId)
}

func DeleteArticle(id int, userID uint) error {
	return models.DeleteArticle(id, userID)
}
func GetArticlesByPage(page, pageSize int) ([]models.Article, int64, error) {
	return models.GetArticlesByPage(page, pageSize)
}
func CreateArticlesBatch(items []models.Article) error {
	return models.CreateArtclesBatch(items)
}
func GetMyArticles(userID uint) ([]models.Article, error) {
	return models.GetMyArticles(userID)
}
