package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer HTTPServerConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type HTTPServerConf struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(configFile string) (*Config, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := NewConfig()
	err = toml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
