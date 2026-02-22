package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Thresholds Thresholds `yaml:"thresholds"`
	Alerts     Alerts     `yaml:"alerts"`
}

type Thresholds struct {
	CPU    int     `yaml:"cpu"`
	Memory int     `yaml:"memory"`
	Disk   int     `yaml:"disk"`
	Load   float64 `yaml:"load"`
}

type Alerts struct {
	Enabled bool   `yaml:"enabled"`
	Webhook string `yaml:"webhook"`
	Email   string `yaml:"email"`
}

func Load(cfgFile string) *Config {
	v := viper.New()
	v.SetDefault("thresholds.cpu", 80)
	v.SetDefault("thresholds.memory", 85)
	v.SetDefault("thresholds.disk", 90)
	v.SetDefault("thresholds.load", 2.0)
	v.SetDefault("alerts.enabled", false)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, _ := os.UserHomeDir()
		v.AddConfigPath(home)
		v.AddConfigPath(".")
		v.SetConfigName(".salam")
		v.SetConfigType("yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		// Config file not found; use defaults
	}

	var cfg Config
	v.Unmarshal(&cfg)
	return &cfg
}

func Generate() error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".salam.yaml")

	defaultConfig := Config{
		Thresholds: Thresholds{
			CPU:    80,
			Memory: 85,
			Disk:   90,
			Load:   2.0,
		},
		Alerts: Alerts{
			Enabled: false,
			Webhook: "",
			Email:   "",
		},
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Configuration created at %s\n", configPath)
	return nil
}
