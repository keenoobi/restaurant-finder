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

type GeneralConfig struct {
	BatchSize int    `yaml:"batch_size"`
	LogLevel  string `yaml:"log_level"`
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
