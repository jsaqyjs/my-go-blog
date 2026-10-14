package models

import (
	"gorm.io/gorm"
	"my-go-blog/database"
)

type Article struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	Author  string `gorm:"type:varchar(100);not null;default:''" json:"author"`
}

func (Article) TableName() string {
	return "articles"
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
func CreateArticle(title, content, author string) (Article, error) {

	// 创建文章结构体
	NewArticle := Article{
		Title:   title,
		Content: content,
		Author:  author,
	}
	err := database.DB.Create(&NewArticle).Error
	return NewArticle, err

}

// 修改文章
func UpdateArticle(id int, title, content, author string) (Article, error) {
	var article Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return Article{}, err
	}
	//if err := database.DB.Model(&article).Updates(map[string]interface{}{
	//	"title":   title,
	//	"content": content,
	//}).Error; err != nil {
	//	return Article{}, err
	//}
	//article.Title = title
	//article.Content = content

	article.Title = title
	article.Content = content

	if err := database.DB.Save(&article).Error; err != nil {
		return Article{}, err
	}

	result := database.DB.Model(&article).Where("id=?", id).Updates(map[string]interface{}{
		"title":   title,
		"content": content,
		"author":  author,
	})
	// 3. 检查数据库错误
	if result.Error != nil {
		return Article{}, result.Error
	}

	// 4. 如果没有影响记录，再确认文章是否存在
	if result.RowsAffected == 0 {
		var count int64

		if err := database.DB.Model(&Article{}).
			Where("id = ?", id).
			Count(&count).Error; err != nil {
			return Article{}, err
		}

		if count == 0 {
			return Article{}, gorm.ErrRecordNotFound
		}
	}
	return article, nil
}

// 删除文章
func DeleteArticle(id int) error {
	result := database.DB.
		Where("id = ?", id).
		Delete(&Article{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

// 搜索title进行找出文章
func SearchArticles(title string) ([]Article, error) {
	var articles []Article
	err := database.DB.
		Where("title LIKE ?", "%"+title+"%").
		Order("id DESC").
		Find(&articles).Error
	return articles, err
}

// 作者查询文章
func GetArticleByAuthor(author string) ([]Article, error) {
	var articles []Article
	err := database.DB.Where("author = ?", author).Find(&articles).Error
	return articles, err
}

// 查询最新文章
func GetLatestArticle() (Article, error) {
	var article Article
	err := database.DB.Order("id desc").First(&article).Error
	return article, err
}

// 返回文章总数
func CountArticles() (int64, error) {
	var count int64
	err := database.DB.Model(&Article{}).Count(&count).Error
	return count, err
}
