package configs

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`	
	Redis    RedisConfig    `mapstructure:"redis"`
	AI       AIConfig       `mapstructure:"ai"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	JWT  JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AIConfig struct {
	Provider string       `mapstructure:"provider"`
	Deepseek DeepseekConfig `mapstructure:"deepseek"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	Gemini   GeminiConfig   `mapstructure:"gemini"`
	Prompts  PromptConfig   `mapstructure:"prompts"`
}

type DeepseekConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	BaseURL string `mapstructure:"base_url"`
}

type GeminiConfig struct {
	APIKey string `mapstructure:"api_key"`
	Model  string `mapstructure:"model"`
}

type PromptConfig struct {
	Summary   string `mapstructure:"summary"`
	Classify  string `mapstructure:"classify"`
	Sentiment string `mapstructure:"sentiment"`
}
