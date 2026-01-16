package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	RPCURL            string        `mapstructure:"rpc_url"`
	WSURL             string        `mapstructure:"ws_url"`
	ContractAddress   string        `mapstructure:"contract_address"`
	StartBlock        uint64        `mapstructure:"start_block"`
	ChainID           uint64        `mapstructure:"chain_id"`
	ReconnectInterval time.Duration `mapstructure:"reconnect_interval"`
	MaxRetries        int           `mapstructure:"max_retries"`
	NFTMetadata       NFTMetadataConfig `mapstructure:"nft_metadata"`
}

// NFTMetadataConfig NFT 元数据配置
type NFTMetadataConfig struct {
	IPFSGateway      string        `mapstructure:"ipfs_gateway"`       // IPFS 网关地址
	ArweaveGateway   string        `mapstructure:"arweave_gateway"`   // Arweave 网关地址
	HTTPTimeout      time.Duration `mapstructure:"http_timeout"`     // HTTP 请求超时时间
	UpdateTimeout    time.Duration `mapstructure:"update_timeout"`   // 元数据更新超时时间
	UserAgent        string        `mapstructure:"user_agent"`       // HTTP User-Agent
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"` // 过期时间（小时）
	Issuer     string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `mapstructure:"level"` // debug, info, warn, error
}

// Global 全局配置实例
var Global *Config

// Load 加载配置文件
// 根据环境变量 GO_ENV 加载对应的配置文件，默认为 local
func Load() *Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 支持环境变量覆盖配置
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("读取配置文件失败", "error", err, "env", env)
		os.Exit(1)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("解析配置文件失败", "error", err)
		os.Exit(1)
	}

	Global = &cfg
	slog.Info("配置加载成功", "env", env, "config_file", viper.ConfigFileUsed())
	return &cfg
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
