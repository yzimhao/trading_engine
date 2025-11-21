# 开发指南（DEVELOPMENT.md）

本文档面向开发者，说明如何在本地搭建、调试、测试与生成文档的常用步骤。

## 环境与依赖
- Go 版本：以 `go.mod` 为准（本仓库注明 `go 1.23.8`）；建议使用 `gvm` / `asdf` 等管理工具将本地 Go 版本与项目一致。
- Docker & Docker Compose：用于启动依赖服务（Postgres、Redis、RocketMQ 等）。
- 推荐安装工具：`golangci-lint`, `swag`（用于生成 swagger docs）, `mockgen`（如需生成 mock）。可通过 `make install` 安装部分工具。

## 本地快速启动
1. 克隆仓库并进入目录：

```bash
git clone https://github.com/yzimhao/trading_engine.git
cd trading_engine
```

2. 启动依赖服务：

```bash
docker compose up -d
```

注意：本仓库的 `docker-compose.yml` 将容器内部端口映射到宿主机端口：Postgres `5432` 映射为宿主 `15432`，Redis `6379` 映射为宿主 `16379`。
因此在宿主（本机）直接连接数据库或 redis 时，请使用对应的宿主端口（例如 `localhost:15432`、`localhost:16379`）。

3. 运行数据库迁移（开发环境）：

```bash
make migrate-up
```

4. 启动服务：

```bash
make run
# 或
bash scripts/run.sh
```

5. 启动前端示例（可选）：

```bash
cd frontend/example
npm install
npm run dev:h5
```

## 迁移策略（重要）
- 开发：自动迁移（自动创建表与新增字段）可以提高迭代速度。
- 生产：严禁在未审核的情况下直接运行自动迁移。仓库中保留 `migrations/postgres/*.sql` 用于受控变更，并在 CI/CD 中以受审计流程运行。
- 本仓库已调整自动迁移行为：当表已存在时，迁移逻辑仅会为模型中不存在的字段添加列（AddColumn）；不会修改已有字段或变更类型，从而降低生产破坏性。请在发布前确认所有 schema 更改是否有对应的手动 SQL 迁移脚本。

## 测试与本地检查
- 运行测试：`go test ./...`
- 模块/包测试：`go test ./pkg/matching -run TestSomething`
- 静态检查：`gofmt -w .`、`golangci-lint run`、`go vet ./...`。

## 生成接口文档（Swagger）
- 生成命令：

```bash
make docs-gen
```

- 生成后文档位于：`generated/docs`，可将其部署在文档服务器或通过本地静态文件服务预览。

## 生成 Mock（测试时使用）
- 脚本：`scripts/mockgen.sh`。该脚本会根据 `internal` 下的接口生成 mocks 到 `mocks/`。

## 调试技巧
- 使用 `dlv`（Delve）进行断点调试：

```bash
dlv debug ./cmd/main -- --your-flags
```

- 在本地调试时可通过环境变量或 `config` 覆盖 listen/port 等配置（见 `provider.NewViper` 相关实现）。

## 常见命令汇总

```bash
# 启动依赖
docker compose up -d

# 迁移
make migrate-up

# 启动服务
make run

# 生成 docs
make docs-gen

# 运行所有测试
go test ./...

# 运行 lint
golangci-lint run
```

## 回滚与应急
- 如果自动迁移引发问题，请暂停服务并按以下步骤回滚：
  1. 停止服务
  2. 恢复数据库备份（事先应在发布前进行备份）
  3. 本地复现问题并生成手动 migration SQL 修复脚本

## 贡献流程小结
- 新增 schema/字段时：
  - 优先编写手动 migration SQL 并提交到 `migrations/postgres/`；
  - 在 PR 描述中说明是否需要运行迁移脚本与回滚步骤；
  - 在 CI 中运行迁移校验（可选脚本）。

欢迎提 PR 改进本指南或补充项目中的脚本与 CI 流水线。
