package config

import (
	"fmt"
	"time"

	"github.com/vrischmann/envconfig"
)

// Config struct
type (
	Config struct {
		MysqlMaster *Database
	}

	Database struct {
		DSN         string
		ReadTimeout time.Duration `envconfig:"default=1s"`
	}
)

// InitConfig func
func InitConfig(prefix string) (*Config, error) {
	config := &Config{}
	if err := envconfig.InitWithPrefix(config, prefix); err != nil {
		return nil, fmt.Errorf("init config failed: %w", err)
	}

	return config, nil
}
