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
// 评论模块测试
// ========================

// createTestPost 创建测试文章，返回文章ID和token
func createTestPost(t *testing.T) (uint, string) {
	token := registerAndLogin(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	createReq := dto.CreatePostRequest{
		Title:   "测试评论的文章",
		Content: "这是一篇用于测试评论的文章",
	}
	w := performRequest(t, http.MethodPost, "/api/v1/posts", createReq, authHeader)
	resp := parseResponse(t, w)
	var post vo.PostVO
	parseData(t, resp, &post)

	return post.ID, token
}

// TestCommentCreate 测试创建评论
func TestCommentCreate(t *testing.T) {
	postID, token := createTestPost(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	t.Run("创建评论成功", func(t *testing.T) {
		req := dto.CreateCommentRequest{
			Content: "这是一条测试评论",
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodPost, path, req, authHeader)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var comment vo.CommentVO
		parseData(t, resp, &comment)

		if comment.Content != req.Content {
			t.Errorf("期望内容 %s, 实际 %s", req.Content, comment.Content)
		}
		if comment.PostID != postID {
			t.Errorf("期望文章ID %d, 实际 %d", postID, comment.PostID)
		}
		if comment.ID == 0 {
			t.Error("评论ID不应为0")
		}
	})

	t.Run("未认证创建评论", func(t *testing.T) {
		req := dto.CreateCommentRequest{
			Content: "未认证评论",
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodPost, path, req, nil)
		assertStatus(t, w, http.StatusUnauthorized)
	})

	t.Run("参数验证失败-内容为空", func(t *testing.T) {
		req := dto.CreateCommentRequest{
			Content: "",
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodPost, path, req, authHeader)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("文章不存在", func(t *testing.T) {
		req := dto.CreateCommentRequest{
			Content: "对不存在的文章评论",
		}
		w := performRequest(t, http.MethodPost, "/api/v1/posts/99999999/comments", req, authHeader)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeNotFound)
	})

	t.Run("其他用户可以评论", func(t *testing.T) {
		// 其他用户也能评论
		anotherToken := registerAndLogin(t)
		anotherAuthHeader := map[string]string{"Authorization": "Bearer " + anotherToken}

		req := dto.CreateCommentRequest{
			Content: "其他用户的评论",
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodPost, path, req, anotherAuthHeader)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)
	})
}

// TestCommentGetByPostID 测试获取文章评论列表
func TestCommentGetByPostID(t *testing.T) {
	postID, token := createTestPost(t)
	authHeader := map[string]string{"Authorization": "Bearer " + token}

	// 创建几条评论
	for i := 0; i < 3; i++ {
		req := dto.CreateCommentRequest{
			Content: fmt.Sprintf("测试评论 %d", i+1),
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		performRequest(t, http.MethodPost, path, req, authHeader)
	}

	t.Run("获取评论列表成功", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodGet, path, nil, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var commentList vo.CommentListVO
		parseData(t, resp, &commentList)

		if commentList.Total < 3 {
			t.Errorf("期望至少3条评论, 实际 %d", commentList.Total)
		}
		if len(commentList.List) == 0 {
			t.Error("评论列表不应为空")
		}
	})

	t.Run("分页参数", func(t *testing.T) {
		path := fmt.Sprintf("/api/v1/posts/%d/comments?page=1&page_size=2", postID)
		w := performRequest(t, http.MethodGet, path, nil, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var commentList vo.CommentListVO
		parseData(t, resp, &commentList)

		if len(commentList.List) > 2 {
			t.Errorf("分页限制2条, 实际返回 %d 条", len(commentList.List))
		}
	})

	t.Run("文章不存在", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts/99999999/comments", nil, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeNotFound)
	})

	t.Run("无效的文章ID", func(t *testing.T) {
		w := performRequest(t, http.MethodGet, "/api/v1/posts/invalid/comments", nil, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})
}

// TestCommentMultipleUsers 测试多用户评论场景
func TestCommentMultipleUsers(t *testing.T) {
	// 用户A创建文章
	postID, _ := createTestPost(t)

	// 用户B和C评论
	for i := 0; i < 2; i++ {
		userToken := registerAndLogin(t)
		userAuthHeader := map[string]string{"Authorization": "Bearer " + userToken}

		req := dto.CreateCommentRequest{
			Content: fmt.Sprintf("用户%d的评论", i+1),
		}
		path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
		w := performRequest(t, http.MethodPost, path, req, userAuthHeader)
		assertStatus(t, w, http.StatusOK)
	}

	// 验证评论数量
	path := fmt.Sprintf("/api/v1/posts/%d/comments", postID)
	w := performRequest(t, http.MethodGet, path, nil, nil)
	resp := parseResponse(t, w)

	var commentList vo.CommentListVO
	parseData(t, resp, &commentList)

	if commentList.Total < 2 {
		t.Errorf("期望至少2条评论, 实际 %d", commentList.Total)
	}
}
