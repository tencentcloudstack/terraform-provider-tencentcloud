## Context

`tencentcloud_teo_deploy_config_group_version` 是 RESOURCE_KIND_GENERAL 资源，负责在 EdgeOne (TEO) 中发布配置组版本。其生命周期为：

- **Create**：调用 `DeployConfigGroupVersion` 发布版本，随后轮询 `DescribeEnvironments`（service 层 `DescribeTeoEnvironmentsByFilter`，入参 `ZoneId`）等待目标 `EnvId` 的环境状态从 `creating`/`version_deploying` 变为 `running`。复合 ID 形如 `zoneId#envId#recordId`。
- **Read**：从复合 ID 解析 `zoneId`、`envId`、`recordId`，通过 `DescribeTeoDeployConfigVersionHistoryByFilter`（按 `record-id` 过滤）读取部署记录，填充 `record_id`、`deploy_time`、`status`、`message`、`description`、`config_group_version_infos`、`zone_id`、`env_id`。
- **Delete**：no-op（资源代表一次性部署动作，不支持回滚删除）。

当前 Read 仅消费 `DescribeConfigGroupVersionHistory` 的部署记录，未将 `DescribeEnvironments` 返回的环境维度信息写入 state。本次变更在 Read 中补充对 `DescribeEnvironments` 的调用与字段映射，使 state 包含环境实际生效状态。

云 API 数据结构（vendor `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`）已确认：

- `DescribeEnvironmentsResponseParams`：
  - `TotalCount *uint64`
  - `EnvInfos []*EnvInfo`
- `EnvInfo`：
  - `EnvId *string`
  - `EnvType *string`（`production` / `staging`）
  - `Status *string`（`creating` / `running` / `version_deploying`）
  - `Scope []*string`（生产为 `["ALL"]`，测试为节点 IP 列表）
  - `CurrentConfigGroupVersionInfos []*ConfigGroupVersionInfo`
  - `CreateTime *string`
  - `UpdateTime *string`
- `ConfigGroupVersionInfo`：
  - `VersionId`、`VersionNumber`、`SourceVersion`、`GroupType`、`GroupId`、`Description`、`Status`、`CreateTime`（均 `*string`）

注意：现有资源已存在名为 `config_group_version_infos` 的字段，其数据来源于 `DeployRecord.ConfigGroupVersionInfos`（部署记录视角），与本次新增的 `current_config_group_version_infos`（来源于 `EnvInfo.CurrentConfigGroupVersionInfos`，环境当前生效版本视角）语义不同，故不复用、不覆盖。

## Goals / Non-Goals

**Goals:**
- 在 Read 中补充 `DescribeEnvironments` 调用，定位目标环境并映射其字段为 Computed 出参。
- 新增字段：`total_count`、`env_type`、`scope`、`current_config_group_version_infos`（嵌套集合）、`env_create_time`、`env_update_time`。
- 保持向后兼容：仅新增 Computed 字段，不修改既有 schema 与既有字段语义。
- 遵循项目硬约束：调用云 API 使用 retry（`tccommon.ReadRetryTimeout`）+ `tccommon.RetryError`；set 前判 nil；空返回先打日志再 `d.SetId("")`。

**Non-Goals:**
- 不修改 `config_group_version_infos`（部署记录视角）字段。
- 不新增/修改 service 层方法（复用现有 `DescribeTeoEnvironmentsByFilter`）。
- 不调整 Create/Delete 逻辑（Create 已使用同一 service 方法轮询状态，保持不变）。
- 不为新增字段提供用户输入能力（Computed-only）。
- 不修改资源 ID 格式。

## Decisions

### 决策 1：复用现有 service 层方法 `DescribeTeoEnvironmentsByFilter`

现有 `DescribeTeoEnvironmentsByFilter(ctx, param)` 接受 `param["ZoneId"]`，返回 `[]*EnvInfo`（取 `response.Response.EnvInfos`）。其请求 `DescribeEnvironmentsRequest` 仅含 `ZoneId`，符合本次读取需求。

- **选择**：直接在 Read 中调用该 service 方法，入参 `paramMap["ZoneId"] = &zoneId`。
- **备选**：新增按 `EnvId` 过滤的 service 方法。**不采纳**，因为 `DescribeEnvironmentsRequest` 无 `EnvId`/`Filters` 字段（仅 `ZoneId`），无法在 API 侧按环境过滤；且现有方法已满足需求，避免重复代码。
- **环境定位**：在返回的 `EnvInfos` 切片中遍历查找 `*item.EnvId == envId` 的元素。

### 决策 2：新增 Computed 字段 schema 定义

- `total_count`：`schema.TypeInt`，`Computed: true`。（云 API 为 `uint64`，Terraform schema 用 `TypeInt` 即可承载。）
- `env_type`：`schema.TypeString`，`Computed: true`。
- `scope`：`schema.TypeList`，`Elem: schema.TypeString`，`Computed: true`。（`EnvInfo.Scope` 为 `[]*string`，映射为字符串列表。）
- `current_config_group_version_infos`：`schema.TypeSet`，`Computed: true`，`Elem` 为 `schema.Resource`，其 schema 含 `version_id` / `version_number` / `source_version` / `group_type` / `group_id` / `description` / `status` / `create_time`，均为 `schema.TypeString`、`Computed: true`。
  - **使用 TypeSet 而非 TypeList**：与现有 `config_group_version_infos` 一致，且配置组版本集合无序，Set 语义更合适。
- `env_create_time`：`schema.TypeString`，`Computed: true`。
- `env_update_time`：`schema.TypeString`，`Computed: true`。
- **命名**：用 `env_create_time` / `env_update_time` 区分于 `ConfigGroupVersionInfo.CreateTime`（已存在于 `current_config_group_version_infos` 子结构内）及部署记录字段，避免同名冲突。

### 决策 3：Read 中的调用顺序与空值处理

Read 方法在解析 `zoneId`/`envId`/`recordId` 后，先保留现有部署记录读取逻辑，随后追加环境信息读取：

1. 现有逻辑：调用 `DescribeTeoDeployConfigVersionHistoryByFilter` 读取部署记录；若返回空，先打 `[CRUD]` 日志保留 id，再 `d.SetId("")` 返回 nil（既有实现已如此，保持不变，不动该段）。
2. 新增逻辑：构造 `paramMap["ZoneId"] = &zoneId`，用 retry 包裹调用 `DescribeTeoEnvironmentsByFilter`（Read 超时 `tccommon.ReadRetryTimeout`，错误用 `tccommon.RetryError` 包装）。
   - retry 块内仅做 API 调用与返回，不做 set/定位（遵循"retry 块只调用接口"约束）。
   - retry 成功后，在 retry 块外遍历 `respData` 定位 `EnvId == envId` 的 `EnvInfo`。
   - 若未定位到目标环境（含返回空切片）：打 `log.Printf("[CRUD] teo deploy config group version env not found, zone_id=%s env_id=%s", zoneId, envId)`，不报错、不 `d.SetId("")`，直接跳过新字段赋值返回（避免误清 ID；部署记录已读取成功即认为资源存在）。
   - 若定位到：按字段逐个判 nil 后 `d.Set(...)`：
     - `EnvType` → `env_type`
     - `Scope`：转为 `[]interface{}`（元素为 string），`d.Set("scope", ...)`
     - `CurrentConfigGroupVersionInfos`：遍历构造 `[]map[string]interface{}`，每个子字段判 nil 后赋值，`d.Set("current_config_group_version_infos", ...)`
     - `CreateTime` → `env_create_time`
     - `UpdateTime` → `env_update_time`
     - `TotalCount`：取 `response.Response.TotalCount`，但 service 层方法只返回 `EnvInfos` 切片不返回 `TotalCount`。
       - **处理**：`TotalCount` 无法从现有 service 方法获得。由于本次目标是"新增出参"，且 service 方法签名返回 `[]*EnvInfo`，调整 service 层签名影响面大（该方法在 Create 轮询中也被调用）。
       - **取舍**：不修改 service 层签名；改为在 Read 中直接调用 `DescribeEnvironments` client 方法一次（经 `meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeEnvironmentsWithContext`），自行读取 `response.Response.TotalCount` 与 `EnvInfos`。
       - **备选 B（采纳）**：扩展 `DescribeTeoEnvironmentsByFilter` 使其同时返回 `TotalCount`（如改返回值或新增方法）。考虑 Create 轮询只需 `EnvInfos`，为避免改动既有调用方，**新增一个 service 层薄封装** `DescribeTeoEnvironmentsWithTotalCount(ctx, zoneId)` 返回 `(uint64, []*EnvInfo, error)`，专供 Read 使用；Create 轮询不变。
       - **最终选择 B**：保持既有方法不动，新增 service 方法，职责清晰、影响面最小。
   - `TotalCount` 在 retry 块内随 API 调用一并取得（service 方法内部调用并返回），定位环境也放在 retry 块外。

### 决策 4：service 层新增方法

在 `service_tencentcloud_teo.go` 新增：

```go
func (me *TeoService) DescribeTeoEnvironmentsWithTotalCount(ctx context.Context, zoneId string) (totalCount uint64, envInfos []*teov20220901.EnvInfo, errRet error)
```

实现：构造 `NewDescribeEnvironmentsRequest()`，`request.ZoneId = helper.String(zoneId)`，`ratelimit.Check`，调用 `DescribeEnvironments`，校验 `response == nil || response.Response == nil` 后取 `TotalCount` 与 `EnvInfos` 返回。Read 中用 `tccommon.ReadRetryTimeout` + `tccommon.RetryError` 包裹。

## Risks / Trade-offs

- **[Risk] `DescribeEnvironments` 仅支持 `ZoneId` 过滤**：当某 Zone 下环境数量较多时，返回全部环境后本地遍历定位 `envId`，存在轻微性能/带宽开销。→ 影响有限（环境数量通常为 2-3 个：生产+测试），且接口已是现有轮询使用的接口，可接受。
- **[Risk] 新增 service 方法与既有 `DescribeTeoEnvironmentsByFilter` 重复调用同一 API**：存在两份近似封装。→ 权衡：避免修改既有方法签名波及 Create 轮询逻辑；新增方法职责单一、便于维护，重复成本可控。
- **[Trade-off] `total_count` 用 `TypeInt` 承载 `uint64`**：极端大值（>2^31）时可能溢出，但环境总数远小于此，且 schema 层面 `TypeInt` 为通用选择。→ 可接受。
- **[Trade-off] 未定位到目标环境时不报错**：当 API 偶发未返回该环境时，新增字段保持未设置。→ 由于部署记录已成功读取（资源存在性由部署记录判定），环境信息缺失不应导致整体 Read 失败或清空 ID；用户可在下次 plan/refresh 时获得刷新。
- **[Risk] retry 块边界**：需严格保证 retry 块内仅 API 调用，环境定位与 set 在块外。→ 通过代码审查与 specs 场景约束保证。

## Migration Plan

- 仅新增 Computed 字段，无数据迁移。
- 部署：更新 provider 二进制后，既有 state 首次 `terraform refresh`/`plan` 时自动填充新字段。
- 回滚：回退 provider 版本即可，新字段忽略不影响既有配置。
