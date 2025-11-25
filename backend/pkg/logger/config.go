package logger

import (
	"os"
	"path/filepath"
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		Production: false,
		Output: OutputConfig{
			File: FileOutputConfig{
				Enabled:    false,
				Path:       "logs/app.log",
				MaxSize:    100,
				MaxAge:     7,
				MaxBackups: 3,
				Compress:   true,
			},
			Console: ConsoleOutputConfig{
				Enabled: true,
				Format:  "console",
				Color:   true,
			},
		},
		Context: ContextConfig{
			AutoFields: []string{"trace_id", "user_id", "org_id"},
			GlobalFields: map[string]interface{}{
				"service": "echomind",
				"version": "0.9.4",
			},
		},
		Sampling: SamplingConfig{
			Enabled: false,
			Rate:    100,
			Burst:   10,
			Levels:  []Level{DebugLevel, InfoLevel},
		},
		Providers: []ProviderConfig{},
	}
}

// ProductionConfig 返回生产环境配置
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.Production = true
	config.Level = InfoLevel
	config.Output.File.Enabled = true
	config.Output.Console.Format = "json"
	config.Output.Console.Color = false
	return config
}

// DevelopmentConfig 返回开发环境配置
func DevelopmentConfig() *Config {
	config := DefaultConfig()
	config.Level = DebugLevel
	config.Output.Console.Enabled = true
	config.Output.Console.Color = true
	return config
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// 日志级别
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		config.Level = parseLevel(levelStr)
	}

	// 生产模式
	if prod := os.Getenv("LOG_PRODUCTION"); prod == "true" || prod == "1" {
		config.Production = true
	}

	// 文件输出
	if path := os.Getenv("LOG_FILE_PATH"); path != "" {
		config.Output.File.Enabled = true
		config.Output.File.Path = path
	}

	// 控制台输出
	if format := os.Getenv("LOG_CONSOLE_FORMAT"); format != "" {
		config.Output.Console.Format = format
	}

	if color := os.Getenv("LOG_CONSOLE_COLOR"); color == "false" || color == "0" {
		config.Output.Console.Color = false
	}

	// 确保日志目录存在
	if config.Output.File.Enabled {
		if err := os.MkdirAll(filepath.Dir(config.Output.File.Path), 0755); err != nil {
			// 如果创建目录失败，禁用文件输出
			config.Output.File.Enabled = false
		}
	}

	return config
}

// parseLevel 解析日志级别字符串
func parseLevel(s string) Level {
	switch s {
	case "DEBUG", "debug":
		return DebugLevel
	case "INFO", "info":
		return InfoLevel
	case "WARN", "warn":
		return WarnLevel
	case "ERROR", "error":
		return ErrorLevel
	case "FATAL", "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}
