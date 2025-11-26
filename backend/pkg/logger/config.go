package logger

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
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

// LoadConfigFromFile 从YAML文件加载配置
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 处理环境变量替换
	data = expandEnvVars(data)

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	// 确保日志目录存在
	if config.Output.File.Enabled {
		if err := os.MkdirAll(filepath.Dir(config.Output.File.Path), 0755); err != nil {
			// 如果创建目录失败，禁用文件输出
			config.Output.File.Enabled = false
		}
	}

	return config, nil
}

// expandEnvVars 展开YAML中的环境变量 ${VAR:default}
func expandEnvVars(data []byte) []byte {
	// 简单的环境变量替换，处理 ${VAR:default} 格式
	content := string(data)

	// 使用正则表达式查找并替换环境变量
	re := regexp.MustCompile(`\$\{([^}:}]+):([^}]*)\}`)

	result := re.ReplaceAllStringFunc(content, func(match string) string {
		// 提取变量名和默认值
		parts := strings.SplitN(match[2:len(match)-1], ":", 2)
		if len(parts) != 2 {
			return match
		}

		varName := parts[0]
		defaultValue := parts[1]

		if envValue := os.Getenv(varName); envValue != "" {
			return envValue
		}
		return defaultValue
	})

	return []byte(result)
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
