# ================================================================
#  gin-admin-api Makefile
#  使用: make <target>   |   make help 查看所有命令
# ================================================================

APP      := gin-admin-api
BINARY   := bin/$(APP)
GO       := go
ENV      := dev

# 版本信息（从 git 读取，无 git 时用 dev）
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date '+%Y-%m-%d %H:%M:%S')
LDFLAGS    := -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)'

GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: all help \
        run runDev runProd dev \
        build build buildMac build-windows \
        wire wireCheck gen \
        tidy fmt vet lint \
        test test-cover \
        docker-build docker-run docker-stop \
        clean chmod \
        startDev startProd

# ── 默认目标 ────────────────────────────────────────────────────
all: tidy wire build

# ── 帮助 ────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "$(GREEN)gin-admin-api 可用命令$(RESET)"
	@echo ""
	@echo "$(YELLOW)▶ 运行$(RESET)"
	@echo "  make run          直接运行（dev 环境）"
	@echo "  make runDev       dev 环境运行"
	@echo "  make runProd      prod 环境运行"
	@echo "  make dev          air 热重载（需安装 air）"
	@echo ""
	@echo "$(YELLOW)▶ 编译$(RESET)"
	@echo "  make build        编译 Linux amd64 二进制"
	@echo "  make buildMac     编译 macOS amd64 二进制"
	@echo "  make build-windows 编译 Windows amd64 二进制"
	@echo ""
	@echo "$(YELLOW)▶ 代码生成$(RESET)"
	@echo "  make wire         重新生成 wire_gen.go（依赖注入）"
	@echo "  make wireCheck    检查 wire 依赖是否正确（不生成文件）"
	@echo "  make gen          运行 gorm-gen 生成类型安全查询代码"
	@echo "  make gen-prod     使用 prod 配置运行 gorm-gen"
	@echo ""
	@echo "$(YELLOW)▶ 代码质量$(RESET)"
	@echo "  make tidy         整理 go.mod / go.sum"
	@echo "  make fmt          格式化所有 Go 文件"
	@echo "  make vet          静态检查"
	@echo "  make lint         golangci-lint 检查"
	@echo ""
	@echo "$(YELLOW)▶ 测试$(RESET)"
	@echo "  make test         运行所有测试"
	@echo "  make test-cover   生成覆盖率报告"
	@echo ""
	@echo "$(YELLOW)▶ Docker$(RESET)"
	@echo "  make docker-build 构建镜像"
	@echo "  make docker-run   启动容器（含 MySQL）"
	@echo "  make docker-stop  停止容器"
	@echo ""
	@echo "$(YELLOW)▶ 部署$(RESET)"
	@echo "  make startDev     pm2 启动（dev 环境）"
	@echo "  make startProd    pm2 启动（prod 环境）"
	@echo "  make chmod        赋予二进制执行权限"
	@echo "  make clean        清理编译产物"
	@echo ""

# ── 运行 ────────────────────────────────────────────────────────

run: wire
	$(GO) run main.go wire_gen.go -envString dev

runDev: wire
	$(GO) run main.go wire_gen.go -envString dev

runProd: wire
	$(GO) run main.go wire_gen.go -envString prod

## dev: air 热重载（需先安装: go install github.com/air-verse/air@latest）
dev:
	@which air > /dev/null 2>&1 || { \
		echo "$(RED)请先安装 air: go install github.com/air-verse/air@latest$(RESET)"; exit 1; }
	air

# ── 编译 ────────────────────────────────────────────────────────

build: wire
	@echo "$(GREEN)▶ 编译 Linux amd64...$(RESET)"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) main.go wire_gen.go
	@echo "$(GREEN)✔ 输出: $(BINARY)$(RESET)"

buildMac: wire
	@echo "$(GREEN)▶ 编译 macOS amd64...$(RESET)"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(APP)-darwin main.go wire_gen.go
	@echo "$(GREEN)✔ 输出: bin/$(APP)-darwin$(RESET)"

build-mac-arm: wire
	@echo "$(GREEN)▶ 编译 macOS arm64 (Apple Silicon)...$(RESET)"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
		$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(APP)-darwin-arm64 main.go wire_gen.go
	@echo "$(GREEN)✔ 输出: bin/$(APP)-darwin-arm64$(RESET)"

build-windows: wire
	@echo "$(GREEN)▶ 编译 Windows amd64...$(RESET)"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
		$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(APP).exe main.go wire_gen.go
	@echo "$(GREEN)✔ 输出: bin/$(APP).exe$(RESET)"

chmod:
	chmod 777 $(BINARY)

# ── Wire 依赖注入 ────────────────────────────────────────────────

## wire: 重新生成 wire_gen.go
## 安装: go install github.com/google/wire/cmd/wire@latest
wire:
	@which wire > /dev/null 2>&1 || { \
		echo "$(YELLOW)wire 未安装，正在安装...$(RESET)"; \
		$(GO) install github.com/google/wire/cmd/wire@latest; }
	@echo "$(GREEN)▶ 生成 wire_gen.go...$(RESET)"
	wire
	@echo "$(GREEN)✔ wire_gen.go 已更新$(RESET)"

## wireCheck: 检查依赖注入是否正确（不生成文件）
wireCheck:
	@which wire > /dev/null 2>&1 || $(GO) install github.com/google/wire/cmd/wire@latest
	wire check
	@echo "$(GREEN)✔ wire 依赖检查通过$(RESET)"

# ── gorm-gen 代码生成 ────────────────────────────────────────────

## gen: 使用 dev 配置生成 GORM 查询代码
gen:
	@echo "$(GREEN)▶ 运行 gorm-gen（dev 环境）...$(RESET)"
	$(GO) run ./internal/query/generator.go application.dev.yml
	@echo "$(GREEN)✔ dao/ 和 model/ 已更新$(RESET)"

## gen-prod: 使用 prod 配置生成 GORM 查询代码
gen-prod:
	@echo "$(GREEN)▶ 运行 gorm-gen（prod 环境）...$(RESET)"
	$(GO) run ./internal/query/generator.go application.prod.yml
	@echo "$(GREEN)✔ dao/ 和 model/ 已更新$(RESET)"

# ── 依赖管理 ─────────────────────────────────────────────────────

tidy:
	@echo "$(GREEN)▶ go mod tidy...$(RESET)"
	$(GO) mod tidy

# ── 代码质量 ─────────────────────────────────────────────────────

fmt:
	@echo "$(GREEN)▶ 格式化代码...$(RESET)"
	$(GO) fmt ./...

vet:
	@echo "$(GREEN)▶ go vet...$(RESET)"
	$(GO) vet ./...

## lint: 需要先安装 golangci-lint
## 安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	@which golangci-lint > /dev/null 2>&1 || { \
		echo "$(RED)请先安装: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; exit 1; }
	golangci-lint run ./...

# ── 测试 ─────────────────────────────────────────────────────────

test:
	@echo "$(GREEN)▶ 运行测试...$(RESET)"
	$(GO) test -v ./...

test-cover:
	@echo "$(GREEN)▶ 生成覆盖率报告...$(RESET)"
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✔ 报告: coverage.html$(RESET)"

# ── Docker ───────────────────────────────────────────────────────

docker-build:
	@echo "$(GREEN)▶ 构建 Docker 镜像...$(RESET)"
	docker build -t $(APP):$(VERSION) -t $(APP):latest .
	@echo "$(GREEN)✔ 镜像: $(APP):$(VERSION)$(RESET)"

docker-run:
	@echo "$(GREEN)▶ 启动容器（含 MySQL）...$(RESET)"
	docker compose up -d
	@echo "$(GREEN)✔ 访问: http://localhost:8080$(RESET)"

docker-stop:
	docker compose down

# ── 部署（pm2）───────────────────────────────────────────────────

startDev:
	pm2 start $(BINARY) -o ./out.log -e ./error.log \
		--log-date-format="YYYY-MM-DD HH:mm Z"

startProd:
	pm2 start $(BINARY) -o ./out.log -e ./error.log \
		--log-date-format="YYYY-MM-DD HH:mm Z" -- -envString prod

# ── 清理 ─────────────────────────────────────────────────────────

clean:
	@echo "$(GREEN)▶ 清理...$(RESET)"
	@rm -rf bin/ coverage.out coverage.html
	$(GO) clean -cache
	@echo "$(GREEN)✔ 清理完成$(RESET)"
