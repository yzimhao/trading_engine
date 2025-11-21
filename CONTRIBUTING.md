# 贡献指南（CONTRIBUTING.md）

感谢你有兴趣为本项目贡献！本文件说明提交 PR、代码风格、测试与本地开发相关的主要流程。

1. 提交前准备
   - Fork 本仓库并在本地创建分支：`git checkout -b feature/短描述`。
   - 分支命名约定：`feature/`, `fix/`, `chore/`, `hotfix/` + `短描述`。
   - 在提交前务必运行测试与静态检查（见下面的“本地开发”章节）。

2. 代码风格与 lint
   - 使用 `gofmt`（或 `go fmt`）格式化代码：`gofmt -w .`
   - 使用 `golangci-lint` 进行静态检查：`golangci-lint run`。将 lint 报告中的高优先级问题修复后再提交。
   - 日志请使用 `go.uber.org/zap`，保持结构化日志；不要用 fmt.Println 打印运行时日志。

3. 本地测试
   - 运行单元测试：`go test ./...`
   - 推荐在变更关键逻辑（撮合、结算、订单）时添加测试用例并确保 CI 通过。

4. 数据库迁移（开发 & CI）
   - 本项目包含自动迁移逻辑（仅用于在开发/测试环境下方便创建表或新增字段）。
   - 生产环境请勿启用自动迁移或请谨慎审核迁移 SQL；建议通过审计过的 migration 文件（`migrations/postgres/*.sql`）并按 CI/CD 审批流程执行。
   - 本地执行迁移示例：

```bash
docker compose up -d      # 启动依赖（Postgres, Redis 等）
make migrate-up           # 运行迁移（开发环境）
```

5. 接口文档（Swagger）
   - 本仓库使用 `swag` 生成 Swagger 文档。
   - 生成命令：`make docs-gen`（内部会调用 `swag init -g cmd/main/main.go --parseDependency --parseInternal -o ./generated/docs`）。
   - 在变更公开 API 时，请更新注释并重新生成文档。

6. 生成 Mock（如需）
   - 项目包含脚本 `scripts/mockgen.sh`，用于生成 mock 文件（用于测试）。

7. 提交与 PR 要点
   - 提交信息请简洁描述变更，若关联 issue 请在 PR 描述中写明。
   - PR 模板应包含：变更概述、如何本地复现、测试覆盖情况、是否需要数据库迁移、回滚建议。

8. 代码审查关注点
   - 检查业务逻辑边界与错误处理（特别是资金/撮合/结算模块）。
   - 检查事务、并发与锁的使用是否正确。
   - 审查与外部依赖（消息、DB、缓存）的契约与超时/重试策略。

9. 联系方式
   - 有疑问请在 issue 中讨论或联系维护者（仓库主页的作者信息）。

谢谢你的贡献！
