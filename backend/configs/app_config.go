package configs

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	AI        AIConfig        `mapstructure:"ai"`
	Security  SecurityConfig  `mapstructure:"security"`
	Worker    WorkerConfig    `mapstructure:"worker"`    // Worker configuration
	Telemetry TelemetryConfig `mapstructure:"telemetry"` // OpenTelemetry configuration
}

type ServerConfig struct {
	Port        string    `mapstructure:"port"`
	Environment string    `mapstructure:"environment"` // "development" | "production"
	JWT         JWTConfig `mapstructure:"jwt"`
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

type WorkerConfig struct {
	Concurrency int `mapstructure:"concurrency"` // Number of concurrent workers
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
	ChunkSize      int                       `mapstructure:"chunk_size"` // Max tokens per chunk for RAG processing
}

type ServiceRoute struct {
	Chat      string `mapstructure:"chat"`
	Embedding string `mapstructure:"embedding"`
}

type ProviderConfig struct {
	Protocol string           `mapstructure:"protocol"` // "openai" | "gemini"
	Settings ProviderSettings `mapstructure:"settings"`
}

type ProviderSettings map[string]interface{}

type PromptConfig struct {
	Summary    string `mapstructure:"summary"`
	Classify   string `mapstructure:"classify"`
	Sentiment  string `mapstructure:"sentiment"`
	DraftReply string `mapstructure:"draft_reply"`
}

// TelemetryConfig defines OpenTelemetry configuration
type TelemetryConfig struct {
	Enabled        bool                    `mapstructure:"enabled"`
	ServiceName    string                  `mapstructure:"service_name"`
	ServiceVersion string                  `mapstructure:"service_version"`
	Environment    string                  `mapstructure:"environment"`
	Exporter       TelemetryExporterConfig `mapstructure:"exporter"`
	Sampling       TelemetrySamplingConfig `mapstructure:"sampling"`
	Metrics        TelemetryMetricsConfig  `mapstructure:"metrics"`
}

type TelemetryExporterConfig struct {
	Type    string                 `mapstructure:"type"` // "console", "file", "otlp"
	Console TelemetryConsoleConfig `mapstructure:"console"`
	File    TelemetryFileConfig    `mapstructure:"file"`
	OTLP    TelemetryOTLPConfig    `mapstructure:"otlp"`
}

type TelemetryConsoleConfig struct {
	EnableColor bool `mapstructure:"enable_color"`
	PrettyPrint bool `mapstructure:"pretty_print"`
}

type TelemetryFileConfig struct {
	TracesPath  string `mapstructure:"traces_path"`
	MetricsPath string `mapstructure:"metrics_path"`
}

type TelemetryOTLPConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`
	Timeout  string `mapstructure:"timeout"`
}

type TelemetrySamplingConfig struct {
	Type  string  `mapstructure:"type"` // "always_on", "always_off", "traceidratio"
	Ratio float64 `mapstructure:"ratio"`
}

type TelemetryMetricsConfig struct {
	ExportInterval string `mapstructure:"export_interval"`
	ExportTimeout  string `mapstructure:"export_timeout"`
}
