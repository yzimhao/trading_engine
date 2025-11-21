# 规范草案 — 单体交易所 Demo

**目标**
- 交付一个单体（monolith）的完整交易所 Demo，包含：用户账户（模拟充值/提现）、资金账本、REST 与 WebSocket API、撮合引擎（限价/市价）、成交持久化、订单管理、简单管理/监控接口。该 Demo 用于展示核心交易流程、资金一致性保障与 E2E 场景验证。

**范围与边界（MVP）**
- 包含：
  - 用户注册/认证（简化 JWT）
  - 账户资产：充值（模拟）、提现（模拟）、转账（内部）
  - 下单：限价、市场单（按数量/按金额）
  - 撮合：单价优先、时间优先，支持撮合并输出成交事件
  - 成交持久化到数据库（trade 表）与资金流水（assets_logs）
  - 订单生命周期：新建、部分成交、完全成交、撤单（含市价成交后自动撤销剩余）
  - API：REST（管理/查询/下单）与 WebSocket（撮合/成交推送）
  - 简易管理面板接口（查看订单/成交/账户快照）
- 不包含（MVP 外）：外部清算对接、KYC、法币网关、复杂风控策略、高级订单类型（止盈止损）、分布式消息中间件依赖（可选用内置队列模拟）

**假设**
- 部署环境为单机 PostgreSQL + 单服务进程（无分布式部署需求）。
- 性能目标为中等（示例：目标能处理 500–2,000 订单/秒的演示负载，依环境不同而异），首要关注正确性与一致性。
- 使用 shopspring/decimal 做小数运算以保证金额精度。

**参与者（Actors）**
- 终端用户（Trader）：通过 REST 或 WebSocket 下单、查询订单、提现/充值（模拟）。
- 管理员（Admin）：查询系统状态、强制取消订单、查看账目快照。
- 撮合引擎（Engine）：接收订单、维护买/卖队列、输出成交与移除事件。
- 持久化消费者（Persistor）：接收成交事件并写入 DB，同时更新账户资产与流水（同进程同步或保证 Ack）。

**关键用例与用户场景**
1. 用户充值 → 下单 → 部分/全部成交 → 查询成交与余额
   - 验证点：充值后可下单，成交后资产与流水正确记录。
2. 用户下市价按金额买入 → 多笔对手单撮合 → 剩余自动撤单
   - 验证点：成交记录顺序与撤单幂等性；不出现双扣/重复记账。
3. 重放请求（相同 transId）确保幂等：重复提交充值/提现请求不重复影响余额
   - 验证点：DB 唯一索引 + repo 层幂等校验生效。
4. WebSocket 实时推送成交给客户端
   - 验证点：客户端能接收并正确展示成交数据，断线重连后可重取快照。

**功能需求（可测试、按优先级 P0/P1/P2）**
- P0（MVP 必须）
  - FR-001: 用户认证（JWT） — 可登录并获取 token（验收：登录接口返回 token，受保护接口需 token）
  - FR-002: 账户充值（模拟） — 提交 `transId, amount, symbol`，保证幂等（验收：重复 transId 不重复记账）
  - FR-003: 账户提现（模拟） — 提交 `transId, amount, symbol`，保证幂等与余额检查（验收：余额不足拒绝）
  - FR-004: 下单（限价/市价） — 支持买/卖，返回 orderId（验收：下单成功，order 存库）
  - FR-005: 撮合与成交输出 — 引擎能正确撮合并产生 trade 事件（验收：trade 写库，订单状态变化）
  - FR-006: 成交持久化与资金流水更新 — trade 写入且资产 ledger 更新（验收：数据库内 trade 与 assets_logs 一致）
  - FR-007: 交易幂等/去重 — `transId` 去重（充值/提现/重要流水）
  - FR-008: 市价订单自动撤单逻辑 — 市价成交结束后自动发布剩余撤单（验收：撤单事件被消费并更新订单状态）

- P1（重要，非 MVP 强制但强烈要求）
  - FR-009: WebSocket 实时推送成交/订单变更
  - FR-010: 订单查询、用户未成交列表、成交历史接口
  - FR-011: 基本管理接口（查看系统快照、强制取消）
  - FR-012: 单机高可用配置（日志、指标）

- P2（改进/后续）
  - FR-013: 充值/提现人工审核流
  - FR-014: 高级风控、订单簿持久化快照

**成功标准（可测量、技术无关）**
- SC-01: 幂等性 — 向充值/提现接口重复提交相同 `transId`，系统不会改变最终账户余额（100% 幂等）。
- SC-02: 资金一致性 — 在 1000 个并发下单与充值混合的 E2E 场景中，结算后账户总资产与流水一致（差异 0）。
- SC-03: 功能可用性 — REST API 与 WS 接口能在 95% 时间内响应性在 200ms 内（演示环境目标）。
- SC-04: 可靠性 — 在模拟延迟与重放场景下，系统能保证成交不会引发重复记账（幂等 + 去重生效）。

**关键实体（最小集合）**
- User { id, username, ... }
- Asset { user_id, symbol, total_balance, avail_balance, freeze_balance }
- Order { id, user_id, symbol, side, price, quantity, filled_qty, status, type, created_at }
- Trade { id, buy_order_id, sell_order_id, price, quantity, timestamp }
- AssetLog { id, user_id, symbol, trans_id, change_type, amount, before, after, created_at }
- Freeze { id, user_id, symbol, amount, trans_id, status }

**数据与接口草案（示例）
- POST /api/v1/account/deposit
  - body: { "transId": "string", "symbol": "USDT", "amount": "123.45" }
  - resp: { "status": "ok", "balance": "..." }

- POST /api/v1/orders
  - body: { "symbol":"BTC_USDT", "side":"buy|sell", "type":"limit|market_by_qty|market_by_amount", "price":"...","quantity":"..." }
  - resp: { "orderId":"...", "status":"accepted" }

- WS /ws (auth with token)
  - subscribe to symbol: { "subscribe": ["ticker:BTC_USDT"] }
  - trade push: { "type":"trade", "symbol":"...", "price":"...", "quantity":"...", "tradeId":"..." }

**一致性策略（建议，MVP 实现细则）**
- 资金变更（充值/提现/资产流水）采用数据库事务并配合 `trans_id` 唯一索引保证幂等：在写入 assets_logs 前尝试插入 trans 记录或直接以 assets_logs.trans_id 建唯一索引，冲突则视为重复请求并返回 OK。此策略适用于单机部署。
- 成交持久化与订单状态更新：引擎撮合在内存生成 trade 事件 → 同步调用持久化函数（在同进程内执行）将 trade 写入 DB 并在同一事务内更新资产流水与订单状态（保证强一致性）。若写入失败，撮合结果应重试或记录失败并暴露给运维。该设计优先保证正确性，牺牲部分吞吐量（符合 demo 目标）。

**测试场景（必须）**
- TS-1: 幂等测试：对充值接口重复发送相同 `transId`，验证余额不变。
- TS-2: 并发下单与撮合：同时产生 1000 个订单，验证最终订单状态及 trade/asset consistency。
- TS-3: 市价按金额买入测试：用户下市价按金额买入，多个卖单被逐一撮合，剩余自动撤单。
- TS-4: WS 推送与断线重连：订阅并接收 trade 推送，断线后重连并能获取最近快照。

**可交付物（MVP）**
- `specs/001-full-exchange/spec.md`（本文件）
- 基础 DB migration（用户、资产、订单、trade、assets_logs、freeze）
- 实现代码：REST API、WS、撮合引擎、持久化逻辑
- 自动化测试：单元测试 + 若干 E2E 脚本（演示流水）
- 部署说明与运行脚本（本地 Docker compose）

**后续步骤（短期路线 3–6 周）**
1. 完成本地运行环境与 DB migration（1–2 天）
2. 实现持久层幂等（transId 唯一）并添加测试（2–4 天）
3. 实现撮合引擎与同步持久化（4–7 天）
4. 实现 API 与 WS（3–5 天）
5. 集成测试与演示资料（2–4 天）

**开放问题（需要明确）**
- OQ-1: 演示是否需要支持多币对并行吞吐的压力测试？（影响性能目标）
- OQ-2: 是否需要在 MVP 中保留现有代码结构并复用部分实现，还是允许从新模块重写？

---

以上为首版规范草案。我会把该文件写入 `specs/001-full-exchange/spec.md`（已完成），并继续推进下一步：收集非功能需求与细化用户旅程（`todo` 已更新）。

请确认：
- 是否接受此 MVP 范围与一致性策略（强一致性，单进程同步持久化）？
- 对开放问题 OQ-1 / OQ-2 的偏好（简单回答即可）？

确认后我会把功能需求拆成更细的任务并更新 `tasks.md`，或直接开始实现你首选的任务（例如 `T001`~`T003` 或 `撮合实现`）。