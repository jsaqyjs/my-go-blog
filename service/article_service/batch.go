package article_service

import (
	"context"
	"fmt"
	"log/slog"
	"my-go-blog/database"
	"my-go-blog/models"
	"sync"
)

type BatchResult struct {
	ArticleID uint   `json:"article_id"`
	Title     string `json:"title,omitempty"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// Worker 需要接收的任务
type articleJob struct {
	Index     int
	ArticleID uint
}

func BatchProcessArticles(
	ctx context.Context,
	articleIDs []uint,
	userID uint,
) []BatchResult {
	wg := sync.WaitGroup{}
	workPool := 3
	// 用于分发任务
	jobs := make(chan articleJob)

	// 预先分配结果切片
	results := make([]BatchResult, len(articleIDs))

	for i := 0; i < workPool; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// articleJob{index, article}:
			for job := range jobs {
				slog.Info(
					"Worker 开始处理文章",
					"worker_id", workerID,
					"article_id", job.ArticleID,
				)
				var article models.Article
				err := database.DB.Where("id = ? AND user_id = ?", job.ArticleID, userID).First(&article).Error
				result := BatchResult{
					ArticleID: job.ArticleID,
				}
				if err != nil {
					result.Error = err.Error()
				} else {
					result.Success = true
					result.Title = article.Title
				}

				results[job.Index] = result
				slog.Info(
					"Worker 处理完成",
					"worker_id", workerID,
					"article_id", job.ArticleID,
				)
			}
			// 保留原来的 GORM 查询与结果赋值代码

		}(i + 1)
	}

dispath:
	for index, article := range articleIDs {
		select {
		case <-ctx.Done():
			break dispath
		case jobs <- articleJob{index, article}:
		}
	}
	close(jobs)
	wg.Wait()
	// 对没有被处理的任务填充结果 因为上面可能导致任务没有完成
	for i, id := range articleIDs {
		if results[i].ArticleID == 0 {
			results[i] = BatchResult{
				ArticleID: id,
				Error:     fmt.Sprintf("任务未完成: %v", ctx.Err()),
			}
		}
	}
	return results
}
