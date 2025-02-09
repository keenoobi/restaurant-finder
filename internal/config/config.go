package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	CSV           CSVConfig           `yaml:"csv"`
	General       GeneralConfig       `yaml:"general"`
	WebApp        WebAppConfig        `yaml:"web"`
	JWT           JWTConfig           `yaml:"jwt"`
}

type ElasticsearchConfig struct {
	Address string                 `yaml:"address"`
	Index   string                 `yaml:"index"`
	Mapping map[string]interface{} `yaml:"mapping"`
}

type CSVConfig struct {
	FilePath  string `yaml:"file_path"`
	Delimiter string `yaml:"delimiter"`
}

type WebAppConfig struct {
	Port     string `yaml:"port"`
	HTMLPage string `yaml:"index_page"`
}

type GeneralConfig struct {
	BatchSize int    `yaml:"batch_size"`
	LogLevel  string `yaml:"log_level"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int64  `yaml:"expiration"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}

	return config, nil
}
