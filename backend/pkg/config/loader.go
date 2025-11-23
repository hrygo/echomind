package config

import (
	"fmt"
	"strings"

	"github.com/hrygo/echomind/configs"
	"github.com/spf13/viper"
)

func Load(path string) (*configs.Config, error) {
	vip := viper.New()
	vip.SetConfigFile(path)
	vip.AddConfigPath(".")
	vip.AutomaticEnv()
	vip.SetEnvPrefix("ECHOMIND")
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg configs.Config
	if err := vip.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Helper: If using viper elsewhere, we might want to return the viper instance too, 
	// but mapping to struct is cleaner.
	// For now, we rely on the struct.
	
	// Inject global Viper fallback if needed by other packages (legacy support)
	// viper.Reset()
	// viper.MergeConfigMap(vip.AllSettings())

	return &cfg, nil
}
