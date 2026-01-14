package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/vo"
	"go_blog/internal/pkg/response"
)

// ========================
// 用户模块测试
// ========================

// 生成唯一用户名（避免测试冲突）
func uniqueUsername() string {
	return fmt.Sprintf("testuser_%d", time.Now().UnixNano())
}

func uniqueEmail() string {
	return fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
}

// TestUserRegister 测试用户注册
func TestUserRegister(t *testing.T) {
	t.Run("注册成功", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: uniqueUsername(),
			Password: "password123",
			Email:    uniqueEmail(),
		}

		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var user vo.UserVO
		parseData(t, resp, &user)

		if user.Username != req.Username {
			t.Errorf("期望用户名 %s, 实际 %s", req.Username, user.Username)
		}
		if user.Email != req.Email {
			t.Errorf("期望邮箱 %s, 实际 %s", req.Email, user.Email)
		}
		if user.ID == 0 {
			t.Error("用户ID不应为0")
		}
	})

	t.Run("用户名已存在", func(t *testing.T) {
		username := uniqueUsername()
		// 第一次注册
		req := dto.RegisterRequest{
			Username: username,
			Password: "password123",
			Email:    uniqueEmail(),
		}
		performRequest(t, http.MethodPost, "/api/v1/register", req, nil)

		// 第二次注册（相同用户名）
		req.Email = uniqueEmail() // 不同邮箱
		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		email := uniqueEmail()
		// 第一次注册
		req := dto.RegisterRequest{
			Username: uniqueUsername(),
			Password: "password123",
			Email:    email,
		}
		performRequest(t, http.MethodPost, "/api/v1/register", req, nil)

		// 第二次注册（相同邮箱）
		req.Username = uniqueUsername() // 不同用户名
		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("参数验证失败-用户名太短", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: "ab", // 少于3个字符
			Password: "password123",
			Email:    uniqueEmail(),
		}

		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("参数验证失败-密码太短", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: uniqueUsername(),
			Password: "12345", // 少于6个字符
			Email:    uniqueEmail(),
		}

		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})

	t.Run("参数验证失败-邮箱格式错误", func(t *testing.T) {
		req := dto.RegisterRequest{
			Username: uniqueUsername(),
			Password: "password123",
			Email:    "invalid-email",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/register", req, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})
}

// TestUserLogin 测试用户登录
func TestUserLogin(t *testing.T) {
	// 先注册一个用户
	username := uniqueUsername()
	password := "password123"
	registerReq := dto.RegisterRequest{
		Username: username,
		Password: password,
		Email:    uniqueEmail(),
	}
	performRequest(t, http.MethodPost, "/api/v1/register", registerReq, nil)

	t.Run("登录成功", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: username,
			Password: password,
		}

		w := performRequest(t, http.MethodPost, "/api/v1/login", req, nil)
		assertStatus(t, w, http.StatusOK)

		resp := parseResponse(t, w)
		assertSuccess(t, resp)

		var loginVO vo.LoginVO
		parseData(t, resp, &loginVO)

		if loginVO.Token == "" {
			t.Error("Token不应为空")
		}
		if loginVO.User.Username != username {
			t.Errorf("期望用户名 %s, 实际 %s", username, loginVO.User.Username)
		}
	})

	t.Run("用户名不存在", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: "nonexistent_user",
			Password: "password123",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/login", req, nil)
		assertStatus(t, w, http.StatusUnauthorized)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeUnauthorized)
	})

	t.Run("密码错误", func(t *testing.T) {
		req := dto.LoginRequest{
			Username: username,
			Password: "wrong_password",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/login", req, nil)
		assertStatus(t, w, http.StatusUnauthorized)

		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeUnauthorized)
	})

	t.Run("参数验证失败-缺少用户名", func(t *testing.T) {
		req := dto.LoginRequest{
			Password: "password123",
		}

		w := performRequest(t, http.MethodPost, "/api/v1/login", req, nil)
		resp := parseResponse(t, w)
		assertCode(t, resp, response.CodeBadRequest)
	})
}

// registerAndLogin 注册并登录，返回 token
func registerAndLogin(t *testing.T) string {
	username := uniqueUsername()
	password := "password123"

	// 注册
	registerReq := dto.RegisterRequest{
		Username: username,
		Password: password,
		Email:    uniqueEmail(),
	}
	performRequest(t, http.MethodPost, "/api/v1/register", registerReq, nil)

	// 登录
	loginReq := dto.LoginRequest{
		Username: username,
		Password: password,
	}
	w := performRequest(t, http.MethodPost, "/api/v1/login", loginReq, nil)
	resp := parseResponse(t, w)

	var loginVO vo.LoginVO
	parseData(t, resp, &loginVO)

	return loginVO.Token
}
