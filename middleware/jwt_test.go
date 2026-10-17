package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"my-go-blog/utils"

	"github.com/gin-gonic/gin"
)

func TestJWTAuthMiddleware(t *testing.T) {

	// 使用仅供测试的 JWT 密钥
	t.Setenv(
		"JWT_SECRET",
		"test-only-secret-1234567890-abcdefghijklmnopqrstuvwxyz",
	)

	token, err := utils.GenerateToken(5)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantUserID uint
	}{
		{
			name:       "缺少Token",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "非法Token",
			authHeader: "Bearer invalid.token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "正确Token",
			authHeader: "Bearer " + token,
			wantStatus: http.StatusOK,
			wantUserID: 5,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			r := gin.New()

			r.GET(
				"/protected",
				JWTAuthMiddleware(),
				func(c *gin.Context) {

					userID, ok := GetCurrentUserID(c)

					if !ok {
						c.JSON(500, gin.H{
							"message": "无法获取用户ID",
						})
						return
					}

					c.JSON(200, gin.H{
						"user_id": userID,
					})
				},
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/protected",
				nil,
			)

			if tt.authHeader != "" {
				req.Header.Set(
					"Authorization",
					tt.authHeader,
				)
			}

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Fatalf(
					"期望状态码%d，实际%d，响应%s",
					tt.wantStatus,
					w.Code,
					w.Body.String(),
				)
			}

			// 成功时进一步检查用户 ID
			if tt.wantStatus == http.StatusOK {
				var body struct {
					UserID uint `json:"user_id"`
				}

				if err := json.Unmarshal(
					w.Body.Bytes(), &body,
				); err != nil {
					t.Fatal(err)
				}

				if body.UserID != tt.wantUserID {
					t.Errorf(
						"期望用户ID=%d，实际%d",
						tt.wantUserID,
						body.UserID,
					)
				}
			}
		})
	}
}
