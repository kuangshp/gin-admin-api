package initialize

import (
	"log"

	"github.com/swaggo/swag/gen"
)

const (
	SwaggerJSONPath = "/swagger.json"
	SwaggerUIPath   = "/swagger/index.html"
)

// GenerateSwaggerDocs generates swagger files before the HTTP server starts.
func GenerateSwaggerDocs() error {
	return gen.New().Build(&gen.Config{
		Debugger:        log.Default(),
		SearchDir:       ".",
		OutputDir:       "docs",
		OutputTypes:     []string{"go", "json", "yaml"},
		MainAPIFile:     "main.go",
		ParseDependency: 3,
		ParseInternal:   true,
		ParseGoList:     false,
	})
}
