package config

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

// LoadEnv 向上递归找 .env 并加载
func LoadEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	// fallback 到默认行为（当前目录）
	return godotenv.Load()
}
