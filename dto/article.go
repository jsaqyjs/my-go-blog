package dto

type ArticleRequest struct {
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
	Author  string `json:"author" binding:"required"`
}
