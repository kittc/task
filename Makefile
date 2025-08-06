.PHONY: help install dev build clean test start-backend start-frontend

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  install          - 安装后端和前端依赖"
	@echo "  dev              - 启动开发环境（并行启动后端和前端）"
	@echo "  build            - 构建项目"
	@echo "  clean            - 清理构建文件"
	@echo "  test             - 运行测试"
	@echo "  start-backend    - 启动后端服务"
	@echo "  start-frontend   - 启动前端服务"

# 安装依赖
install:
	@echo "安装Go依赖..."
	go mod tidy
	@echo "安装Node.js依赖..."
	npm install

# 启动开发环境
dev:
	@echo "启动开发环境..."
	@$(MAKE) -j2 start-backend start-frontend

# 启动后端服务
start-backend:
	@echo "启动后端服务..."
	go run cmd/server/main.go

# 启动前端服务
start-frontend:
	@echo "启动前端服务..."
	npm run dev

# 构建项目
build: build-backend build-frontend

# 构建后端
build-backend:
	@echo "构建后端..."
	CGO_ENABLED=1 go build -o bin/server cmd/server/main.go

# 构建前端
build-frontend:
	@echo "构建前端..."
	npm run build

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf dist/
	rm -rf node_modules/
	go clean

# 运行测试
test:
	@echo "运行Go测试..."
	go test ./...
	@echo "运行前端测试..."
	npm test

# 生成API文档
docs:
	@echo "生成API文档..."
	# 这里可以添加生成API文档的命令

# 数据库迁移
migrate:
	@echo "运行数据库迁移..."
	go run cmd/server/main.go migrate

# 创建生产环境构建
prod-build: clean
	@echo "创建生产环境构建..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/server cmd/server/main.go
	npm run build

# 启动生产环境
prod-start:
	@echo "启动生产环境..."
	GIN_MODE=release ./bin/server