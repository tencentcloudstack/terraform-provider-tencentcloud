## Context

TencentCloud MongoDB 提供审计服务（Audit Service），允许用户对实例开启数据库审计日志。当前 Terraform Provider 已有 MongoDB 相关资源（实例、备份、SSL 等），但缺少审计服务管理能力。

云 API 提供了完整的 CRUD 接口：
- `OpenAuditService` - 开通审计（异步）
- `DescribeAuditConfig` - 查询审计配置
- `ModifyAuditService` - 修改审计配置
- `CloseAuditService` - 关闭审计（异步）

审计服务绑定到 MongoDB 实例，一个实例只能有一个审计服务配置。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_mongodb_audit_service` 资源的完整 CRUD 生命周期
- 支持全审计和规则审计两种模式
- 支持审计日志保存时长配置
- 正确处理异步操作（开通/关闭）的轮询等待
- 提供单元测试（gomonkey mock 方式）

**Non-Goals:**
- 不实现审计日志查询数据源
- 不实现审计规则的独立资源管理
- 不修改现有 MongoDB 资源

## Decisions

### 1. 资源 ID 使用 instance_id

**决策**: 使用 `instance_id` 作为资源 ID（`d.SetId(instanceId)`）

**理由**: 审计服务与实例一对一绑定，`instance_id` 是唯一标识。无需复合 ID。

**替代方案**: 无其他合理选择，API 本身以 instance_id 为唯一标识。

### 2. instance_id 设置为 ForceNew

**决策**: `instance_id` 字段标记为 `ForceNew: true`

**理由**: 审计服务绑定到特定实例，无法将审计服务从一个实例迁移到另一个实例。变更 instance_id 需要销毁重建。

### 3. 异步操作轮询策略

**决策**: Create（OpenAuditService）和 Delete（CloseAuditService）调用后，通过 `DescribeAuditConfig` 轮询 `IsOpening`/`IsClosing` 字段，直到值为 `"false"` 表示操作完成。

**理由**: 这两个接口是异步操作，需要等待后端完成。DescribeAuditConfig 返回的 IsOpening/IsClosing 字段可以判断操作是否仍在进行中。

**替代方案**: 使用固定 sleep 等待 — 不可靠，可能等待不足或过长。

### 4. LogExpireDay 类型统一为 int

**决策**: Terraform Schema 中 `log_expire_day` 使用 `TypeInt`

**理由**: OpenAuditService 中该字段为 `*uint64`，ModifyAuditService 和 DescribeAuditConfig 中为 `*int64`。Terraform Schema 统一使用 int 类型，在调用不同 API 时进行类型转换。

### 5. rule_filters 使用 TypeList 嵌套结构

**决策**: `rule_filters` 使用 `TypeList` + `Elem: &schema.Resource{}`，子字段包含 `type`（string）、`compare`（string）、`value`（list of string）

**理由**: 对应云 API 的 `[]*LogFilter` 结构，每个 LogFilter 包含 Type、Compare、Value 三个字段。

### 6. 单元测试使用 gomonkey mock

**决策**: 使用 gomonkey 对云 API 客户端方法进行 mock，不使用 Terraform 验收测试套件

**理由**: 新增资源要求使用 gomonkey mock 方式进行单元测试，只验证业务逻辑正确性。

## Risks / Trade-offs

- **[异步操作超时]** → 在 Schema 中声明 Timeouts（Create/Delete），轮询时使用 context deadline 控制超时
- **[DescribeAuditConfig 在审计未开通时的行为]** → Read 方法中如果返回空或错误表示资源不存在，需要调用 `d.SetId("")` 移除 state
- **[rule_filters 在 AuditAll=true 时的处理]** → 当 audit_all 为 true 时，rule_filters 应为空或忽略；Read 时 DescribeAuditConfig 不返回 rule_filters，需要在 state 中保持用户配置或清空
