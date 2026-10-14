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
