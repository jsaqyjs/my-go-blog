package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {

	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	// 构造模拟 GET 请求
	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 让 Gin 处理请求
	r.ServeHTTP(w, req)

	// 检查状态码
	if w.Code != http.StatusOK {
		t.Fatalf("期望200，实际%d", w.Code)
	}

	// 检查响应体
	expected := `{"message":"ok"}`

	if w.Body.String() != expected {
		t.Errorf(
			"期望%s，实际%s",
			expected,
			w.Body.String(),
		)
	}
}
