package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env       string `yaml:"env" env-default:"default"`
	SourceDir string `yaml:"sourceDir"`
	LogsDir   string `yaml:"logsDir"`
	PackPromo int    `yaml:"packPromo"`
	DbConfig  `yaml:"db_config"`
	LogConfig `yaml:"log_config"`
}

type DbConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host" env-default:"127.0.0.1"`
	Port     string `yaml:"port" env-default:"3306"`
	DbName   string `yaml:"db_name"`
}

type LogConfig struct {
	DefaultLogFile string `yaml:"defaultLogFile"`
}

func (c *DbConfig) GormDns() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Password, c.Host, c.Port, c.DbName)
}

func MustLoad(rootPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(fmt.Sprintf("%s/config/config.yaml", rootPath), &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/config/env.yaml", rootPath)); err == nil {
		_ = cleanenv.ReadConfig(fmt.Sprintf("%s/config/env.yaml", rootPath), &cfg)
	}

	_ = cleanenv.ReadConfig(fmt.Sprintf("%s/config/%s.yaml", rootPath, cfg.Env), &cfg)

	return &cfg
}
