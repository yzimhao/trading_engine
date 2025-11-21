# 项目深度分析 — trading_engine

日期: 2025-11-21

简短说明：本文档为对当前代码库的深度静态分析与整理，聚焦架构、模块职责、运行与开发指引、关键发现与改进建议，供维护与下一步开发决策参考。

````markdown
# 项目深度分析 — trading_engine (更新)

日期: 2025-11-21

简短说明：基于先前静态分析，并结合仓库中已实现的若干修复（涉及数据库迁移策略、撮合引擎行为、并发优化、Docker 与本地端口/卷映射、以及取消挂单端到端验证尝试），本文档对项目现状、关键改动、已验证行为与下一步推荐做汇总更新。

---

## 一、总体回顾与当前状态
- 项目仍为单体式后端服务，语言/框架与之前一致（Go + Gin + Fx + GORM + Postgres + Redis）。
- 本次实现并验证的重要修复点：
	- 数据库自动迁移逻辑从一次性 `AutoMigrate(...)` 调用，改为按模型逐表处理：新建表或仅为模型新增缺失列（避免修改已存在列）。实现文件：`internal/persistence/database/auto_migrate.go`（新增 `addMissingColumns` 辅助逻辑）。
	- 撮合引擎行为修正（关键业务修复）：`pkg/matching` 的 `Engine.AddItem` 改为“先撮合，剩余再入簿”（active-match-then-insert），避免此前的“先入簿再撮合”带来的价格穿透与订单快照不一致问题。
	- 并发与通知改善：增加 `resultNotify` / `removeNotify` 通道缓冲（从很小的缓冲改为例如 1024），并用异步 `emitTradeResult` / `emitRemoveResult` helper 发送通知，避免在持锁时阻塞 channel 发送导致的死锁/性能问题。
	- 本地开发环境调整：`docker-compose.yml` 的宿主卷映射改为 `./docker-data/...`，同时修改 provider 的默认端口以匹配 Docker 映射（Postgres -> `15432`, Redis -> `16379`）。文件：`docker-compose.yml`，`internal/di/provider/gorm.provider.go`，`internal/di/provider/redis.provider.go`。
	- 前端取消挂单链路测试与验证：找到前端调用 `GET /api/v1/order/cancel`（`frontend/example`），后端 `internal/modules/base/order` 中的 `cancel` 会构建 `EventNotifyCancelOrder` 并调用 `produce.Publish(...)` 发布到 `notify_order_cancel` 主题，但该接口受 JWT 中间件保护，未认证请求会被中间件拦截并返回 Unauthorized，从而不会执行 publish（我在本机用 `curl` 测试得到 `Unauthorized` 并确认 Redis 列表为空）。

## 二、主要改动（按文件/模块）
- `internal/persistence/database/auto_migrate.go`
	- 用逐模型方式替换直接 `AutoMigrate`，新增 `addMissingColumns`：解析模型 schema 并仅为缺失列调用 `Migrator().AddColumn`，跳过没有 DB 列名的字段并改进错误信息。

- `pkg/matching`（多个文件：`engine.go`, `limit_order.go`, `market_order.go` 等）
	- `Engine.AddItem` 改为主动撮合 limit 订单，只有剩余才放入 orderbook。
	- 增大 `resultNotify` 与 `removeNotify` 缓冲，新增 `emitTradeResult` / `emitRemoveResult` 异步发送 helper，避免在持锁时直接写 channel。
	- 测试：已在 `pkg/matching` 运行并通过本地测试（已执行 `go test`，10 个测试通过，race detector 运行未报告竞态）。

- `internal/modules/base/order/order.go`
	- `cancel` 路由实现：生成 `models_types.EventNotifyCancelOrder` 并用 `o.produce.Publish(c, models_types.TOPIC_NOTIFY_ORDER_CANCEL, body)` 发布；注意：路由受 JWT 中间件保护，必须先认证。

- `internal/di/provider/produce.provider.go`
	- `Publish` 使用 Redis `LPUSH` 模拟消息队列（开发用），`Consume.Subscribe` 使用 `BRPOP` 订阅。

- `docker-compose.yml`
	- 卷映射统一到 `./docker-data/*`，并映射 Postgres 主机端口 `15432`、Redis 主机端口 `16379`（便于本地与容器端口对齐）。

## 三、验证与观察（我已在本机执行的步骤与结果）
- 单元测试
	- 已添加并运行验证 `EventNotifyCancelOrder` JSON 格式的单元测试：`internal/modules/base/order/order_cancel_test.go`（测试通过）。

- 启动依赖服务
	- 我已使用仓库根的 `docker-compose.yml` 启动 `postgres`, `redis`, `rocketmq` 相关容器（映射到 `./docker-data`）。

- 启动服务与 E2E 尝试
	- 我在本地用 `nohup make run > app.log 2>&1 &` 启动了应用并调用 `GET /api/v1/order/cancel?symbol=btcusdt&order_id=B123`。
	- 结果：HTTP 200（Gin 响应包装），但响应体为 `{"code":1000,"msg":"Unauthorized"}`；应用日志记录 `cookie token is empty`，说明 JWT 中间件拒绝了未认证请求，导致 `o.produce.Publish` 没被执行，Redis `notify_order_cancel` 列表为空。
	- 我已把后台进程 PID 写入 `./app.pid`，并把日志写在 `./app.log`，随后按你的指示停止了后台进程。

## 四、问题定位与结论
- 数据库迁移
	- 之前直接 `AutoMigrate` 在某些 struct 字段与 DB 映射不明确时会生成不正确 SQL（例如空列名），已通过使用模型 schema 的 `DBName` 字段、跳过无映射字段并逐列添加的实现降低风险。仍建议在生产中禁用自动迁移并使用显式 SQL migration 流程（已记录）。

- 取消挂单 E2E
	- 前端发起取消请求前必须认证（JWT）。目前在 demo 前端或测试脚本中调用 `GET /api/v1/order/cancel` 需要携带 token（cookie 或 Authorization），否则中间件会提前返回，publish 不会发生。

- 撮合引擎
	- 主动优先撮合的变更修复了价格穿透与订单查找不一致的问题（理论上解决了“挂单先入簿再撮合导致穿透/查不到订单”的缺陷）。并发改进减少了在高并发下 channel 阻塞导致的异常风险。

## 五、已完成的变更清单（快速目录）
- `internal/persistence/database/auto_migrate.go` — 逐模型迁移 + `addMissingColumns`。
- `pkg/matching/engine.go`, `pkg/matching/limit_order.go`, `pkg/matching/market_order.go` — AddItem 行为修改、channel 缓冲与 emit helper。
- `internal/di/provider/gorm.provider.go`、`internal/di/provider/redis.provider.go` — 本地端口默认值调整以匹配 `docker-compose.yml`。
- `docker-compose.yml` — 卷与端口映射调整到 `./docker-data` 与 host ports `15432/16379`。
- `internal/modules/base/order/order_cancel_test.go` — 新增单元测试用于验证取消消息格式。

## 六、后续建议与优先行动项（短期）
1. 生产迁移策略：在 `internal/persistence/database` 中提供开关（env）以区分 dev/prod；在生产禁用自动新增列的行为，改用可审计的 SQL migration。优先级：高。
2. 端到端取消流程验证：用两种方式完成验证：
	 - 在测试环境中通过登录接口生成 JWT 并以该 token 调用 `/api/v1/order/cancel`，确认 Redis 有推送并观察 `matching` 模块消费（推荐做法）。
	 - 或在短期内用 `LPUSH notify_order_cancel <payload>` 验证 publish 层格式（快速但不能验证消费）。优先级：中。
3. 增加集成测试与 CI：将 `matching` 的关键集成用例（撮合场景、取消订单）加入 CI，并在 PR 中运行带 `-race` 的测试。优先级：中-高。
4. 文档与示例：在 `frontend/example` 中补充示例登录流程（或在 README 中说明调用 `/api/v1/order/cancel` 需认证），避免开发者误以为该接口可匿名调用。优先级：低。

## 七、下一步（我可以代劳）
- 如果你同意，我可以：
	1. 用项目的登录接口自动获取 JWT 并完成一次完整的取消 E2E 测试（包括检查 Redis、matching 日志与 process 结果），并把结果和关键日志片段贴回；或
	2. 只在 Redis 上模拟 `LPUSH notify_order_cancel <payload>` 并展示 LRANGE 输出（更快，但不验证消费）；或
	3. 在 `internal/persistence/database` 中增加 `AUTOMIGRATE_DEV=true` 之类的开关，明确区分 dev/prod 行为（需要改代码并走 CI）。

请回复你希望我执行的下一步（1/2/3）或者给出其它指示。

---

附：常用命令（项目根）

```bash
# 启动依赖（本地）
docker compose up -d

# 启动服务（本地开发）
nohup make run > app.log 2>&1 &

# 停止后台服务（示例）
kill $(cat app.pid) || true

# 检查 Redis 列表
docker compose exec -T redis redis-cli LRANGE notify_order_cancel 0 -1

# 运行匹配包测试
go test ./pkg/matching -race -v
```

文件来源（部分）：`cmd/main/main.go`, `internal/di/di.go`, `internal/modules/module.go`, `internal/persistence/database/module.go`, `pkg/matching`, `internal/modules/base/order`。

````

