package initialize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/swaggo/swag/gen"
)

const (
	SwaggerJSONPath = "/swagger.json"
	SwaggerUIPath   = "/swagger/index.html"
)

func SwaggerJSONFile() (string, error) {
	root, err := findProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "docs", "swagger.json"), nil
}

// GenerateSwaggerDocs generates swagger files before the HTTP server starts.
func GenerateSwaggerDocs() error {
	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	return gen.New().Build(&gen.Config{
		Debugger:        log.Default(),
		SearchDir:       root,
		OutputDir:       filepath.Join(root, "docs"),
		OutputTypes:     []string{"go", "json", "yaml"},
		MainAPIFile:     "main.go",
		ParseDependency: 3,
		ParseInternal:   true,
		ParseGoList:     false,
	})
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err = os.Stat(filepath.Join(dir, "main.go")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("未找到项目根目录 main.go")
		}
		dir = parent
	}
}
