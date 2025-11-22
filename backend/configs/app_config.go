package configs

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`	
	Redis    RedisConfig    `mapstructure:"redis"`
	AI       AIConfig       `mapstructure:"ai"`
	Security SecurityConfig `mapstructure:"security"` // New security config
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	JWT  JWTConfig `mapstructure:"jwt"`
}

type SecurityConfig struct {
	EncryptionKey string `mapstructure:"encryption_key"`
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
	ActiveServices ServiceRoute              `mapstructure:"active_services"`
	Providers      map[string]ProviderConfig `mapstructure:"providers"`
	Prompts        PromptConfig              `mapstructure:"prompts"`
}

type ServiceRoute struct {
	Chat      string `mapstructure:"chat"`
	Embedding string `mapstructure:"embedding"`
}

type ProviderConfig struct {
	Protocol string                 `mapstructure:"protocol"` // "openai" | "gemini"
	Settings map[string]interface{} `mapstructure:"settings"`
}

type PromptConfig struct {
	Summary    string `mapstructure:"summary"`
	Classify   string `mapstructure:"classify"`
	Sentiment  string `mapstructure:"sentiment"`
	DraftReply string `mapstructure:"draft_reply"`
}
