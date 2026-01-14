package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_blog/internal/app/router"
	"go_blog/internal/pkg/config"
	"go_blog/internal/pkg/database"
	"go_blog/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var testRouter *gin.Engine

// TestMain 测试入口，初始化测试环境
func TestMain(m *testing.M) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 加载测试配置
	cfg := config.Load()

	// 初始化数据库
	database.NewMySQL(cfg)

	// 初始化路由
	testRouter = router.SetupRouter(gin.TestMode)

	// 运行测试
	m.Run()
}

// ========================
// 测试辅助函数
// ========================

// APIResponse 通用响应结构
type APIResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// performRequest 执行 HTTP 请求
func performRequest(t *testing.T, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("序列化请求体失败: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

// parseResponse 解析响应
func parseResponse(t *testing.T, w *httptest.ResponseRecorder) *APIResponse {
	var resp APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v, body: %s", err, w.Body.String())
	}
	return &resp
}

// assertStatus 断言 HTTP 状态码
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	if w.Code != expected {
		t.Errorf("期望状态码 %d, 实际 %d, body: %s", expected, w.Code, w.Body.String())
	}
}

// assertCode 断言业务状态码
func assertCode(t *testing.T, resp *APIResponse, expected int) {
	if resp.Code != expected {
		t.Errorf("期望业务码 %d, 实际 %d, msg: %s", expected, resp.Code, resp.Msg)
	}
}

// assertSuccess 断言成功响应
func assertSuccess(t *testing.T, resp *APIResponse) {
	if resp.Code != response.CodeSuccess {
		t.Errorf("期望成功响应, 实际 code=%d, msg=%s", resp.Code, resp.Msg)
	}
}

// parseData 解析响应数据
func parseData(t *testing.T, resp *APIResponse, v interface{}) {
	if err := json.Unmarshal(resp.Data, v); err != nil {
		t.Fatalf("解析响应数据失败: %v", err)
	}
}
