package test

import (
	"fmt"
	"net/http"
	"testing"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/vo"
	"go_blog/internal/pkg/response"
)

// ========================
// 文章模块测试
// ========================

// TestPostCreate 测试创建文章
func TestPostCreate(t *testing.T) {
	token := registerAndLogin(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	t.Run("创建文章成功", func(t *testing.T) {
		req := dto.CreatePostRequest{
			Title:   "测试文章标题",
			Content: "测试文章内容，这是一篇测试文章。",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/posts", req, authHeader)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var post vo.PostVO
		parseData(t, resp, &post)

		if post.Title != req.Title {
			t.Errorf("期望标题 %s, 实际 %s", req.Title, post.Title)
		}
		if post.Content != req.Content {
			t.Errorf("期望内容 %s, 实际 %s", req.Content, post.Content)
		}
		if post.ID == 0 {
			t.Error("文章ID不应为0")
		}
	})

	t.Run("未认证创建文章", func(t *testing.T) {
		req := dto.CreatePostRequest{
			Title:   "测试文章",
			Content: "测试内容",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/posts", req, nil)
		assertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("参数验证失败-标题为空", func(t *testing.T) {
		req := dto.CreatePostRequest{
			Title:   "",
			Content: "测试内容",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/posts", req, authHeader)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("参数验证失败-内容为空", func(t *testing.T) {
		req := dto.CreatePostRequest{
			Title:   "测试标题",
			Content: "",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/posts", req, authHeader)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})
}

// TestPostGetList 测试获取文章列表
func TestPostGetList(t *testing.T) {
	t.Run("获取文章列表", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts", nil, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var postList vo.PostListVO
		parseData(t, resp, &postList)

		// 列表不为 nil
		if postList.List == nil {
			t.Error("文章列表不应为 nil")
		}
	})

	t.Run("分页参数", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts?page=1&page_size=5", nil, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)
	})
}

// TestPostGetByID 测试获取文章详情
func TestPostGetByID(t *testing.T) {
	// 先创建一篇文章
	token := registerAndLogin(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	createReq := dto.CreatePostRequest{
		Title:   "测试获取详情",
		Content: "测试内容",
	}
	w := performRequest(t, http.MethodPost, "/api/v1/posts", createReq, authHeader)
	resp := parseResponse(t, w)
	var createdPost vo.PostVO
	parseData(t, resp, &createdPost)

	t.Run("获取文章详情成功", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodGet, path, nil, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var post vo.PostVO
		parseData(t, resp, &post)

		if post.ID != createdPost.ID {
			t.Errorf("期望ID %d, 实际 %d", createdPost.ID, post.ID)
		}
	})

	t.Run("文章不存在", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts/99999999", nil, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeNotFound)
	})

	t.Run("无效的文章ID", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts/invalid", nil, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})
}

// TestPostUpdate 测试更新文章
func TestPostUpdate(t *testing.T) {
	// 创建用户和文章
	token := registerAndLogin(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	createReq := dto.CreatePostRequest{
		Title:   "原始标题",
		Content: "原始内容",
	}
	w := performRequest(t, http.MethodPost, "/api/v1/posts", createReq, authHeader)
	resp := parseResponse(t, w)
	var createdPost vo.PostVO
	parseData(t, resp, &createdPost)

	t.Run("更新文章成功", func(t *testing.T) {
		updateReq := dto.UpdatePostRequest{
			Title:   "更新后的标题",
			Content: "更新后的内容",
		}
		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodPut, path, updateReq, authHeader)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var post vo.PostVO
		parseData(t, resp, &post)

		if post.Title != updateReq.Title {
			t.Errorf("期望标题 %s, 实际 %s", updateReq.Title, post.Title)
		}
	})

	t.Run("未认证更新文章", func(t *testing.T) {
		updateReq := dto.UpdatePostRequest{
			Title: "新标题",
		}
		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodPut, path, updateReq, nil)
		assertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("非作者更新文章", func(t *testing.T) {
		// 用另一个用户尝试更新
		anotherToken := registerAndLogin(t)
		anotherAuthHeader := map[string]string{"Authorization": "Bearer " + anotherToken}

		updateReq := dto.UpdatePostRequest{
			Title: "非法更新",
		}
		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodPut, path, updateReq, anotherAuthHeader)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeForbidden)
	})
}

// TestPostDelete 测试删除文章
func TestPostDelete(t *testing.T) {
	// 创建用户和文章
	token := registerAndLogin(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	createReq := dto.CreatePostRequest{
		Title:   "待删除的文章",
		Content: "这篇文章将被删除",
	}
	w := performRequest(t, http.MethodPost, "/api/v1/posts", createReq, authHeader)
	resp := parseResponse(t, w)
	var createdPost vo.PostVO
	parseData(t, resp, &createdPost)

	t.Run("非作者删除文章", func(t *testing.T) {
		anotherToken := registerAndLogin(t)
		anotherAuthHeader := map[string]string{"Authorization": "Bearer " + anotherToken}

		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodDelete, path, nil, anotherAuthHeader)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeForbidden)
	})

	t.Run("删除文章成功", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/posts/%d", createdPost.ID)
		w := performRequest(t, http.MethodDelete, path, nil, authHeader)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		// 验证文章已被删除
		w = performRequest(t, http.MethodGet, path, nil, nil)
		resp = parseResponse(t, w)
		assertCode(t, resp, response.CodeNotFound)
	})

	t.Run("删除不存在的文章", func(t *testing.T) {
		w := performRequest(t, http.MethodDelete, "/api/v1/posts/99999999", nil, authHeader)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeNotFound)
	})
}
