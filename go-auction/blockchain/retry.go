package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries        int           // 最大重试次数
	InitialBackoff    time.Duration // 初始退避时间
	MaxBackoff        time.Duration // 最大退避时间
	BackoffMultiplier float64       // 退避倍数
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxRetries:        3,
	InitialBackoff:    1 * time.Second,
	MaxBackoff:        60 * time.Second,
	BackoffMultiplier: 2.0,
}

// RetryWithBackoff 带指数退避的重试机制
func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
	backoff := config.InitialBackoff
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}

		// 执行函数
		err := fn()
		if err == nil {
			if attempt > 0 {
				slog.Debug("重试成功", "attempt", attempt)
			}
			return nil
		}

		lastErr = err

		// 如果还有重试机会，等待后重试
		if attempt < config.MaxRetries {
			slog.Warn("操作失败，准备重试", "attempt", attempt+1, "max_retries", config.MaxRetries, "backoff", backoff, "error", err)

			// 等待退避时间
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during backoff: %w", ctx.Err())
			case <-time.After(backoff):
				// 继续重试
			}

			// 计算下次退避时间（指数增长）
			backoff = time.Duration(float64(backoff) * config.BackoffMultiplier)
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}
	}

	return fmt.Errorf("max retries (%d) exceeded, last error: %w", config.MaxRetries, lastErr)
}

// RetryWithBackoffDefault 使用默认配置的重试
func RetryWithBackoffDefault(ctx context.Context, fn func() error) error {
	return RetryWithBackoff(ctx, DefaultRetryConfig, fn)
}
