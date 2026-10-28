## Context

Terraform Provider for TencentCloud 的 PostgreSQL 产品线已存在 `tencentcloud_postgresql_readonly_instance` 资源，但它依赖 `memory` + `cpu` 间接推导 `SpecCode`，且 update/delete 逻辑较为复杂，使用了已不推荐的做法（如 `d.Partial`、直接 `d.SetId("")` 而无现场日志等）。

本次新增 `tencentcloud_postgresql_readonly_instance_v2` 作为 RESOURCE_KIND_GENERAL 通用资源，直接将云 API `CreateReadOnlyDBInstance` 的入参 1:1 暴露为 Terraform schema（包含 `spec_code`、`instance_count` 等新参数），并遵循最新的代码生成规范（retry 错误处理、NonRetryableError 空值检查、现场日志、immutableArgs 等）。

### 云 API 能力确认（vendor 中已验证）
- `CreateReadOnlyDBInstance`：创建只读实例，返回 `DealNames`、`BillId`、`DBInstanceIdSet`、`BillingParameters`。
- `DescribeDBInstanceAttribute`：根据 `DBInstanceId` 查询实例详情，返回 `DBInstance` 结构体（含 zone、vpc、subnet、status、tags 等丰富字段）。
- `IsolateDBInstances`：根据 `DBInstanceIdSet` 隔离实例（删除的前置操作）。

注意：云 API 未提供独立的 "UpdateReadOnlyDBInstance" 接口，实例创建后的属性变更（如 name、storage、project_id 等）需要通过其它 Modify 接口完成。本资源为 GENERAL 资源，需实现 CRUD；对于 Update，鉴于本次接口映射中仅提供了 Create/Describe/Isolate 三个接口，Update 方法将基于 immutableArgs 模式处理：仅 `Id()` 为 ForceNew，其余顶层字段在 update 时若检测到变更则返回 error（CRD-only 资源模式）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_postgresql_readonly_instance_v2` 资源，支持通过 Terraform 创建、查询、隔离（销毁）PostgreSQL 只读实例。
- 将 `CreateReadOnlyDBInstance` 的全部入参暴露为 schema 字段。
- 在 Read 中正确回填 `DescribeDBInstanceAttribute` 返回的字段。
- 遵循最新代码生成规范：retry 包装错误、空值 NonRetryableError、现场日志、set 前 nil 检查。
- 提供基于 gomonkey mock 的单元测试。
- 生成文档与 provider 注册。

**Non-Goals:**
- 不修改既有 `tencentcloud_postgresql_readonly_instance`（v1）资源的行为。
- 不实现复杂的 update 逻辑（不使用额外的 Modify 接口），因为本次接口映射仅包含 Create/Describe/Isolate 三个接口，Update 采用 immutableArgs 模式。
- 不实现 `make doc`（由收尾阶段 tfpacer-finalize 统一执行）。

## Decisions

### 1. 资源 ID 策略
- 使用 `CreateReadOnlyDBInstance` 返回的 `DBInstanceIdSet[0]` 作为资源 ID（单实例场景，`instance_count` 虽可大于 1，但 Terraform 资源只管理一个主实例 ID）。
- 资源支持 import，使用 `schema.ImportStatePassthrough`。

### 2. Schema 字段设计
- 必填（Required, ForceNew）：`zone`、`master_db_instance_id`、`spec_code`、`storage`、`instance_count`、`period`。
- 可选（Optional）：`vpc_id`、`subnet_id`、`instance_charge_type`、`auto_voucher`、`voucher_ids`、`auto_renew_flag`、`project_id`、`activity_id`、`read_only_group_id`、`tag_list`、`security_group_ids`、`need_support_ipv6`、`name`、`db_version`、`dedicated_cluster_id`、`deletion_protection`、`tags`。
- 计算属性（Computed）：`deal_names`、`bill_id`、`db_instance_id_set`、`billing_parameters`、`db_instance_id`（从 ID 回填）。
- `tag_list` 映射为单对象（云 API `Tag` 结构），`tags` 映射为列表；为简化实现，`tag_list` 使用 `TypeList` 单元素结构，`tags` 使用标准 tag map/list。

### 3. Create 流程
1. 构造 `CreateReadOnlyDBInstanceRequest`，填充全部入参。
2. `resource.Retry(tccommon.WriteRetryTimeout)` 调用，失败用 `tccommon.RetryError()` 包装。
3. 检查返回值：`response == nil` 或 `Response == nil` → NonRetryableError。
4. 检查 `DBInstanceIdSet` 是否为空：打印 logId 与 d.Id()，若为空返回 NonRetryableError。
5. 设置 `d.SetId(DBInstanceIdSet[0])`（retry 块外）。
6. 轮询 `DescribeDBInstanceAttribute` 直到实例状态为 running（异步生效）。
7. 设置计算属性 `deal_names`、`bill_id`、`db_instance_id_set`、`billing_parameters`。
8. 处理 `tags`（若提供，通过 tag service 绑定）。

### 4. Read 流程
1. 从 `d.Id()` 获取 instanceId。
2. `resource.Retry(tccommon.ReadRetryTimeout)` 调用 `DescribeDBInstanceAttribute`。
3. 若 `response == nil` 或 `DBInstance == nil`：打印现场日志后 `d.SetId("")`。
4. 按 nil 检查逐个 `d.Set()` 回填字段。
5. 回填 `db_instance_id` = d.Id()。

### 5. Update 流程（immutableArgs 模式）
- 仅 `Id()` 字段为 ForceNew（通过将必填字段设为 ForceNew 实现）。
- Update 方法中将其余顶层字段加入 `immutableArgs` 数组，若检测到变更则返回 error。
- 由于实际可变字段（如 name、project_id 等）需调用额外 Modify 接口，而本次接口映射未提供，故采用 immutable 模式保持简洁。

### 6. Delete 流程
1. 调用 `IsolateDBInstances`（通过 service 层封装），retry 包装错误。
2. 轮询 `DescribeDBInstanceAttribute` 直到状态为 isolated。
3. 注：本次接口映射仅包含 Isolate（隔离），不含 Destroy（销毁）。Delete 方法实现为隔离后即返回，与云 API 能力对齐。

### 7. 单元测试策略
- 使用 gomonkey mock `PostgresqlService` 方法（或直接 mock client 调用）。
- 测试 Create/Read/Delete 的业务逻辑分支（成功、空返回、错误）。
- 不使用 terraform 测试套件。

## Risks / Trade-offs

- [Update 能力受限] → 由于接口映射仅含 Create/Describe/Isolate，Update 采用 immutable 模式，用户如需修改属性需重建资源。这是与云 API 能力对齐的取舍。
- [instance_count > 1] → `instance_count` 允许创建多个实例，但 Terraform 资源只管理 `DBInstanceIdSet[0]` 对应的实例，其余实例不在 state 管理范围。在 schema 描述中说明此限制。
- [tags 与 tag_list 并存] → 云 API 同时提供 `TagList`（单对象，已不推荐）和 `Tags`（列表），两者均暴露为可选字段，由用户选择使用。
- [Delete 仅为隔离] → 与既有 v1 资源（isolate + destroy）不同，v2 的 Delete 仅执行 Isolate，因为接口映射未包含 Destroy 接口。
