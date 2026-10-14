package models

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articles = []Article{
	{
		ID:      1,
		Title:   "Go语言学习",
		Content: "今天开始学习Go后端开发",
	},
	{
		ID:      2,
		Title:   "Gin框架入门",
		Content: "使用Gin实现HTTP接口",
	},
	{
		ID:      3,
		Title:   "MySQL数据库",
		Content: "学习数据库的基本操作",
	},
}

func GetArticles() []Article {
	return articles
}

func GetArticle(id int) (Article, bool) {
	for _, article := range articles {
		if article.ID == id {
			return article, true
		}
	}
	return Article{}, false
}
func CreateArticle(title, content string) Article {
	newID := len(articles) + 1
	// 创建文章结构体
	NewArticle := Article{
		ID:      newID,
		Title:   title,
		Content: content,
	}
	articles = append(articles, NewArticle)
	return NewArticle

}

// 修改文章
func UpdateArticle(id int, title, content string) (Article, bool) {
	for i, article := range articles {
		if article.ID == id {
			articles[i].Content = content
			articles[i].Title = title
			return article, true
		}
	}
	return Article{}, false

}

// 删除文章
func DeleteArticle(id int) bool {
	for i, article := range articles {
		if article.ID == id {
			articles = append(articles[:i], articles[i+1:]...)
			return true
		}
	}
	return false

}
