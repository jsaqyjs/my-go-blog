package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int)

	var wg sync.WaitGroup

	// 创建 3 个 Worker
	for workerID := 1; workerID <= 3; workerID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// 不断从任务队列接收任务
			for articleID := range jobs {
				fmt.Printf(
					"Worker %d 开始处理文章 %d\n",
					id,
					articleID,
				)

				time.Sleep(time.Second)

				fmt.Printf(
					"Worker %d 完成文章 %d\n",
					id,
					articleID,
				)
			}
		}(workerID)
	}

	// 投递 20 个任务
	for articleID := 1; articleID <= 20; articleID++ {
		jobs <- articleID
	}

	// 不再发送新任务
	close(jobs)

	// 等待全部 Worker 完成
	wg.Wait()

	fmt.Println("全部文章处理完成")
}
