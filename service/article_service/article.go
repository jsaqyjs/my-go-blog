package article_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"my-go-blog/database"
	"my-go-blog/models"
	"time"
)

const notFoundCacheValue = "__NOT_FOUND__"

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
func DeleteArticleCache(id int) error {
	return models.DeleteArticleCache(id)
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

func GetArticleCache(id int) (models.Article, error) {
	// Redis Key，例如 article:1
	if id <= 0 {
		return models.Article{}, gorm.ErrRecordNotFound
	}
	key := fmt.Sprintf("article:%d", id)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()
	// 1. 优先从 Redis 读取
	cached, err := database.RedisDB.Get(ctx, key).Result()
	if err == nil {
		// Redis 中缓存的是“文章不存在”
		if cached == notFoundCacheValue {
			return models.Article{}, gorm.ErrRecordNotFound
		}

		var article models.Article

		// 把 Redis 中的 JSON 转换成 Go 结构体
		if json.Unmarshal([]byte(cached), &article) == nil {

			return article, nil
		}
		// 缓存格式异常，继续查询 MySQL
		log.Printf("文章缓存 JSON 解析失败: %s", key)
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("Redis 读取失败: %v", err)
	}
	// 2. 缓存不存在或者读取失败，查询 MySQL
	article, err := models.GetArticleById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			cacheCtx, cacheCancel := context.WithTimeout(
				context.Background(),
				2*time.Second,
			)
			defer cacheCancel()

			cacheErr := database.RedisDB.Set(
				cacheCtx,
				key,
				notFoundCacheValue,
				30*time.Second,
			).Err()

			if cacheErr != nil {
				log.Printf("空值缓存写入失败: %v", cacheErr)
			}
		}
		return models.Article{}, err
	}
	// 3. 将 Article 转换成 JSON
	data, err := json.Marshal(article)
	if err != nil {
		return article, nil
	}

	// 4. 缓存 5 分钟
	cacheCtx, cacheCancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cacheCancel()
	_ = database.RedisDB.Set(
		cacheCtx,
		key,
		data,
		5*time.Minute,
	).Err()
	return article, nil

}
