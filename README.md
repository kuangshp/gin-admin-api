# gin-admin-api

基于 Gin + GORM + Wire 的后台管理 API 项目

## 项目结构

```
.
├── main.go                    # 程序入口
├── wire.go                    # Wire 依赖注入定义
├── wire_gen.go                # Wire 自动生成
├── Makefile                   # 构建脚本
├── Dockerfile                 # Docker 构建
├── docker-compose.yml         # 开发环境 Docker Compose
├── docker-compose.prod.yml    # 生产环境 Docker Compose
├── application.dev.yml        # 开发配置
├── application.prod.yml       # 生产配置
├── go.mod / go.sum            # Go 依赖
│
├── cmd/                       # 可执行命令
│   ├── generator.go           # gorm-gen 代码生成工具
│   └── migrate.go             # 数据库迁移工具
│
├── config/                    # 配置相关
│   └── config.go
│
├── initialize/                # 初始化模块
│   ├── app.go                 # App 入口
│   ├── config.go              # 配置加载
│   ├── db.go                  # 数据库连接
│   ├── logger.go              # 日志初始化
│   ├── redis.go               # Redis 连接
│   ├── router.go              # 路由注册
│   ├── validate.go            # 验证器
│   └── initSql.go             # 数据库初始化（默认账号）
│
├── internal/                  # 内部包
│   ├── api/account/           # 账号 API
│   │   ├── handler.go         # 处理器
│   │   ├── dto/               # 数据传输对象
│   │   └── vo/                # 视图对象
│   ├── middleware/            # 中间件
│   ├── query/                 # 数据访问层
│   │   ├── dao/               # DAO（gorm-gen 生成）
│   │   ├── model/             # 实体模型
│   │   │   ├── entity/        # 数据库实体
│   │   │   └── types/         # 自定义类型
│   │   └── repository/        # 仓储层
│   └── router/                # 路由
│
└── pkg/                       # 公共包
    ├── constants/              # 常量
    ├── enum/                  # 枚举
    ├── utils/                 # 工具函数
    └── validators/            # 自定义验证器
```

## 目录说明

### cmd/
可执行命令目录，包含数据库迁移、代码生成等工具。

### config/
配置相关模块，负责加载 application.yml 配置文件。

### initialize/
应用初始化模块，包含：
- `app.go` - 应用入口，负责启动和优雅关闭
- `config.go` - 配置加载（使用 Viper）
- `db.go` - 数据库连接初始化
- `logger.go` - Zap 日志初始化
- `redis.go` - Redis 连接初始化
- `router.go` - Gin 路由注册
- `validate.go` - Gin 验证器注册
- `initSql.go` - 数据库初始化（默认管理员账号）

### internal/
内部包，不对外暴露。

### pkg/
公共工具包，可被外部项目引用。

## 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+

### 启动项目

```bash
# 1. 配置数据库连接 (application.dev.yml)

# 2. 代码生成（如有表结构变更）
make gen

# 3. 数据库迁移
make migrate

# 4. 运行项目
make run
```

### Makefile 命令

| 命令 | 说明 |
|------|------|
| `make run` | 运行项目（dev 环境） |
| `make build` | 编译 Linux amd64 二进制 |
| `make wire` | 重新生成 wire_gen.go |
| `make gen` | 运行 gorm-gen 生成代码 |
| `make migrate` | 执行数据库迁移 |
| `make tidy` | 整理 go.mod |
| `make test` | 运行测试 |
| `make docker-up` | 启动 Docker 容器 |

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM + gorm-gen
- **依赖注入**: Google Wire
- **配置管理**: Viper
- **日志**: Zap
- **缓存**: Redis
- **验证**: go-playground/validator
