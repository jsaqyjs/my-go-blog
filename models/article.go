package models

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"my-go-blog/database"
	"time"
)

type Article struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	Author  string `gorm:"type:varchar(100);index:idx_author" json:"author"`
	UserID  uint   `gorm:"index;not null" json:"user_id"`
}

var ErrArticleForbidden = errors.New("无权操作此文章")

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
func CreateArticle(title, content, author string, userId uint) (Article, error) {

	// 创建文章结构体
	NewArticle := Article{
		Title:   title,
		Content: content,
		Author:  author,
		UserID:  userId,
	}
	err := database.DB.Create(&NewArticle).Error

	return NewArticle, err

}

// 修改文章
func UpdateArticle(id int, title, content, author string, userID uint) (Article, error) {
	var article Article
	if err := database.DB.First(&article, id).Error; err != nil {
		return Article{}, err
	}

	if article.UserID != userID {
		return Article{}, ErrArticleForbidden
	}
	// 文章所属作者一样 进行修改
	result := database.DB.
		Model(&Article{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"title":   title,
			"content": content,
		})
	if result.Error != nil {
		return Article{}, result.Error
	}
	// 4. 重新查询文章 重新付给article
	err := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&article).Error

	if err != nil {
		return Article{}, err
	}
	return article, nil
}

// 删除文章
func DeleteArticle(id int, userID uint) error {
	var article Article

	// 1. 查询文章
	if err := database.DB.First(&article, id).Error; err != nil {
		return err
	}
	// 2. 判断是否属于当前用户
	if article.UserID != userID {
		return ErrArticleForbidden
	}

	result := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Article{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

// 修改文章删除缓存
func DeleteArticleCache(id int) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	key := fmt.Sprintf("article:%d", id)

	return database.RedisDB.Del(ctx, key).Err()
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

func GetArticlesByPage(page, pageSize int) ([]Article, int64, error) {
	var articles []Article
	var total int64
	err := database.DB.Model(&Article{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * 10
	err = database.DB.Limit(pageSize).Offset(offset).Find(&articles).Error
	return articles, total, err
}

func CreateArtclesBatch(items []Article) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			article := item

			if article.Title == "FAIL" {
				return fmt.Errorf("模拟第二次插入失败")
			}
			if err := tx.Create(&article).Error; err != nil {
				return err
			}
		}
		return nil
	})

}

// 获取当前用户文章
func GetMyArticles(userID uint) ([]Article, error) {
	var articles []Article
	err := database.DB.Model(&Article{}).Where("user_id = ?", userID).
		Order("id DESC").Find(&articles).Error
	return articles, err
}
