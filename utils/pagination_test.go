package utils

import "testing"

func TestCalculateOffset(t *testing.T) {

	tests := []struct {
		name      string
		page      int
		pageSize  int
		want      int
		wantError bool
	}{
		{"第一页", 1, 10, 0, false},
		{"第二页", 2, 10, 10, false},
		{"第三页", 3, 5, 10, false},
		{"页码错误", 0, 10, 0, true},
		{"每页数量错误", 1, 101, 0, true},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got, err := CalculateOffset(
				tt.page,
				tt.pageSize,
			)

			if (err != nil) != tt.wantError {
				t.Fatalf("错误状态不符合预期: %v", err)
			}

			if err != nil {
				return
			}

			if got != tt.want {
				t.Errorf(
					"期望 %d，实际得到 %d",
					tt.want,
					got,
				)
			}
		})
	}
}
