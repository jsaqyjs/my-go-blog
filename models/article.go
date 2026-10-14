package models

import (
	"gorm.io/gorm"
	"my-go-blog/database"
)

type Article struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
}

func GetArticles() ([]Article, error) {
	var articles []Article
	err := database.DB.Order("id ASC").Find(&articles).Error
	return articles, err
}

func GetArticleById(id int) (Article, error) {
	var article Article
	err := database.DB.First(&article, id).Error
	return article, err
}
func CreateArticle(title, content string) (Article, error) {

	// 创建文章结构体
	NewArticle := Article{
		Title:   title,
		Content: content,
	}
	err := database.DB.Create(&NewArticle).Error
	return NewArticle, err

}

// 修改文章
func UpdateArticle(id int, title, content string) (Article, error) {
	var article Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return Article{}, err
	}
	if err := database.DB.Model(&article).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
	}).Error; err != nil {
		return Article{}, err
	}
	article.Title = title
	article.Content = content

	return article, nil
}

// 删除文章
func DeleteArticle(id int) error {
	result := database.DB.Delete(&Article{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}
