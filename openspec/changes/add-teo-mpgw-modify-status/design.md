## Context

`tencentcloud_teo_multi_path_gateway` 目前封装了 TEO Multi-Path Gateway 的 CRUD，其中 Update 仅调用 `ModifyMultiPathGateway`（修改名称、IP、端口），而网关的"启用/停用"由独立 API `ModifyMultiPathGatewayStatus` 负责。现有 `status` 字段为 Computed，仅用于展示状态，无法驱动状态变更。

当前状态：
- `resource_tc_teo_multi_path_gateway.go` 中 schema 已声明 `status` 为 `Computed`
- Update 函数分支仅检查 `gateway_name` / `gateway_ip` / `gateway_port`
- SDK 提供 `teov20220901.NewModifyMultiPathGatewayStatusRequest()`，请求字段：`ZoneId`、`GatewayId`、`GatewayStatus`（取值：`online` / `offline`）

约束：
- Terraform Provider 向后兼容：现有 TF 配置（未写 `status`）必须无 plan drift
- 异步 API：`ModifyMultiPathGatewayStatus` 调用后，实际状态切换需要轮询 `DescribeMultiPathGateways` 确认
- schema 已存在 `status` 字段，只能在 `Computed` 基础上追加 `Optional`（不能改类型或删除）

## Goals / Non-Goals

**Goals:**
- 支持通过 Terraform 配置声明式切换 Multi-Path Gateway 的启用状态
- 保持向后兼容：现有资源升级后不触发 plan 差异
- Update 流程幂等：用户 apply 相同状态不触发 API 调用
- 通过单元测试覆盖 status 变更路径

**Non-Goals:**
- 不改变 Create/Read/Delete 的行为（Create 默认由云端决定初始状态）
- 不对 `gateway_name` / `gateway_ip` / `gateway_port` 的现有 Update 逻辑做任何修改
- 不新增独立的 `tencentcloud_teo_multi_path_gateway_status` 资源
- 不引入新的 Timeouts 子块结构变化（复用资源已有 / 默认超时）

## Decisions

### Decision 1: 在现有资源上扩展 Update，而非新建独立资源

**选择**：复用 `tencentcloud_teo_multi_path_gateway`，把 `status` 改为 `Optional + Computed`，Update 内部多一个分支。

**备选**：新建 `tencentcloud_teo_multi_path_gateway_status` 资源（类似 `confirm_*` 风格的 operation 资源）。

**理由**：
- `status` 语义上属于网关本身属性，单独资源会导致资源数量膨胀且状态分裂（两个资源管同一个 gateway）
- `Optional + Computed` 是 Provider 处理"可选的、未配置时由云端回填"字段的标准模式，与现有 `gateway_port` / `region_id` / `gateway_ip` 一致
- 向后兼容性好：未配置 status 时 `d.GetOk("status")` 返回 `ok=false`，走原有 Computed 路径

### Decision 2: Update 内独立分支，调用 `ModifyMultiPathGatewayStatus`

**选择**：在现有 `needChange` 分支之外，新增一个独立分支：

```go
if d.HasChange("status") {
    if v, ok := d.GetOk("status"); ok {
        // 调用 ModifyMultiPathGatewayStatus
    }
}
```

**备选**：将 status 合并进 `mutableArgs` 一起触发 `ModifyMultiPathGateway`。

**理由**：
- `ModifyMultiPathGateway` API 并不支持 `GatewayStatus` 字段，两个是不同 API
- 分两个独立分支避免条件耦合；任一字段变化都能独立触发对应 API

### Decision 3: 仅在 `d.GetOk("status")` 为 true 时调用 API

**选择**：只有当 `status` 在配置中被显式设置为非空字符串时才调用 `ModifyMultiPathGatewayStatus`。

**理由**：
- 避免首次 import 或旧 state 升级后因 `status` 从 Computed 变为 Optional 而误触发 API
- 用户移除 `status` 字段时，不应将其解读为"关闭网关"，语义上是"交回云端管理"

### Decision 4: 调用后轮询等待状态达成

**选择**：`ModifyMultiPathGatewayStatus` 后通过 `service.DescribeTeoMultiPathGatewayById` 轮询，直到 `respData.Status` 等于期望值，或超时返回错误。复用 `d.Timeout(schema.TimeoutUpdate)`。

**备选**：不等待，立即返回（依赖 Read 下一轮刷新）。

**理由**：
- 与 Create/Update 现有流程一致（Create 也会等待 `online`）
- 避免 apply 结束后立刻 Read 拿到中间态导致 plan drift

### Decision 5: `status` 合法取值约束

**选择**：文档中说明合法取值为 `online` / `offline`；schema 层面**不使用** `ValidateFunc` 限制（与现有字段保持一致，服务端会返回错误）。

**理由**：
- 与当前 provider 的大部分 Optional 字段一致，减少与未来 API 扩展的耦合
- 错误由云 API 返回即可

## Risks / Trade-offs

- **Risk**：已有 state 中 `status` 为 Computed 的空值（或任意中间态），用户首次添加 `status = "online"` 且恰好云端也是 `online` 时，`d.HasChange` 可能因 state 中历史值触发一次幂等调用 → **Mitigation**：`ModifyMultiPathGatewayStatus` 本身对同状态调用幂等（SDK 文档确认），即便触发也无副作用
- **Risk**：轮询超时导致 apply 失败，但实际云端操作已成功 → **Mitigation**：沿用 provider 统一的 `RetryError` 模式，用户可重跑 apply 收敛
- **Trade-off**：把 `status` 做成 Optional+Computed 后，Read 回写与用户配置可能在瞬时过渡态产生短暂不一致 → 可接受，Terraform 会在下一次 refresh 收敛

## Migration Plan

- 新增字段属性为纯加法（Optional 追加），无 state 迁移需求
- 存量资源：Terraform state 中已有 `status` 值，升级后 `terraform plan` 对未在 HCL 配置的资源不会产生 diff（Optional+Computed 行为）
- 文档更新：在 `resource_tc_teo_multi_path_gateway.md` 中补充 status 字段说明和示例
- 回滚：若需要回退，只需将 schema 中 `Optional` 去掉即可（Computed 保留），Update 分支移除；state 值不会丢失

## Open Questions

- 无
