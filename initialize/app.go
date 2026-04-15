package initialize

import (
	"context"
	"fmt"
	"gin-admin-api/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// App 封装 HTTP 服务生命周期
type App struct {
	cfg    *config.ServerConfig
	engine *gin.Engine
	logger *zap.Logger
}

// NewApp Wire Provider
func NewApp(cfg *config.ServerConfig, engine *gin.Engine, logger *zap.Logger) *App {
	return &App{cfg: cfg, engine: engine, logger: logger}
}

// InitSqlData 调用初始化管理员账号
func (a *App) InitSqlData() error {
	return InitAccountDataWithDao()
}

// Run 启动服务，监听系统信号优雅退出
func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      a.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		a.logger.Info("服务启动", zap.String("addr", addr))
		fmt.Printf("✔ 服务已启动: http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("正在优雅关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
