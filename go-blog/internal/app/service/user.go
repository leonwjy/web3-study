package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/model"
	"go_blog/internal/app/repository"
	"go_blog/internal/app/vo"
	"go_blog/internal/pkg/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("用户不存在")
	ErrUsernameExists     = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已被注册")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
)

// UserService 用户服务
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterRequest) (*vo.UserVO, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		slog.Error("检查用户名失败", "error", err)
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// 检查邮箱是否已存在
	exists, err = s.repo.ExistsByEmail(req.Email)
	if err != nil {
		slog.Error("检查邮箱失败", "error", err)
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// 生成盐
	salt, err := generateSalt()
	if err != nil {
		slog.Error("生成盐失败", "error", err)
		return nil, err
	}

	// 使用 bcrypt 加密密码（盐会被 bcrypt 自动处理，这里的 salt 作为额外安全层）
	hashedPassword, err := hashPassword(req.Password, salt)
	if err != nil {
		slog.Error("密码加密失败", "error", err)
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Salt:     salt,
		Email:    req.Email,
	}

	if err := s.repo.Create(user); err != nil {
		slog.Error("创建用户失败", "error", err)
		return nil, err
	}

	slog.Info("用户注册成功", "user_id", user.ID, "username", user.Username)
	return vo.ToUserVO(user), nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*vo.LoginVO, error) {
	// 查找用户
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		slog.Error("查询用户失败", "error", err)
		return nil, err
	}

	// 验证密码
	if !checkPassword(req.Password, user.Salt, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// 生成 JWT Token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		slog.Error("生成Token失败", "error", err)
		return nil, err
	}

	slog.Info("用户登录成功", "user_id", user.ID, "username", user.Username)
	return &vo.LoginVO{
		Token: token,
		User:  *vo.ToUserVO(user),
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*vo.UserVO, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return vo.ToUserVO(user), nil
}

// generateSalt 生成随机盐
func generateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// hashPassword 使用 bcrypt 加密密码
func hashPassword(password, salt string) (string, error) {
	// 将密码和盐拼接后进行 bcrypt 加密
	combined := password + salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkPassword 验证密码
func checkPassword(password, salt, hashedPassword string) bool {
	combined := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combined))
	return err == nil
}
