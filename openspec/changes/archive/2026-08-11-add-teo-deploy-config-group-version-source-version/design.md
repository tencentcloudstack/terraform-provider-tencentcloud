## Context

The `tencentcloud_teo_deploy_config_group_version` resource (file `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.go`) manages the release of EdgeOne config-group versions into an environment. Its `config_group_version_infos` nested block currently exposes these fields: `version_id` (Required), `version_number`, `group_id`, `group_type`, `description`, `status`, and `create_time` (all Computed).

The cloud API struct `teov20220901.ConfigGroupVersionInfo` (in the vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` package) defines an additional field `SourceVersion *string` — the version ID from which a version was derived. This field is returned by `DescribeDeployHistory` at `response.Records[].ConfigGroupVersionInfos[].SourceVersion` (struct `DeployRecord` → `ConfigGroupVersionInfos`), which the resource Read already uses via `TeoService.DescribeTeoDeployConfigVersionHistoryByFilter`.

**约束：**
- 必须保持向后兼容：新增字段只能是 `Computed`，不得改变 Create 行为。
- 现有 Read 函数已遍历 `deployRecord.ConfigGroupVersionInfos` 并逐字段 `set`，新增字段沿用同一遍历逻辑即可。

**关键文件：**
- `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.go:46-86` — `config_group_version_infos` 嵌套 schema
- `tencentcloud/services/eo/resource_tc_teo_deploy_config_group_version.go:281-317` — Read 函数中 `config_group_version_infos` 的回填逻辑
- `vendor/.../teo/v20220901/models.go:2026-2050` — `ConfigGroupVersionInfo` 结构体（含 `SourceVersion`）

## Goals / Non-Goals

**Goals:**
- 在 `config_group_version_infos` 嵌套 schema 中新增 `source_version` 字段（Computed），映射云 API `ConfigGroupVersionInfo.SourceVersion`。
- 在 Read 函数中从 `DescribeDeployHistory` 响应回填 `source_version`，保持 nil 安全（字段为 nil 时不 set）。
- 同步更新资源文档与单元测试。

**Non-Goals:**
- 不修改 Create 流程：`SourceVersion` 仅作为只读展示，不作为 `DeployConfigGroupVersion` 的入参（创建请求 `ConfigGroupVersionInfo` 仅设置 `VersionId`）。
- 不修改 Delete 流程（资源无 Delete API，Delete 为空实现）。
- 不新增任何 ForceNew 或 Required 字段。
- 不引入 `_extension.go` 文件。

## Decisions

### Decision 1: 字段属性为 Computed，不参与 Create
**决策：** `source_version` 设为 `Computed: true`，不设 `Required`/`Optional`/`ForceNew`，Create 函数不读取、不传递该字段。

**理由：**
- 现有同块字段（`version_number`、`group_id`、`group_type`、`description`、`status`、`create_time`）均为 Computed，遵循同一模式。
- 云 API `DeployConfigGroupVersion` 请求中 `ConfigGroupVersionInfo` 仅需 `VersionId`；`SourceVersion` 是版本派生关系的只读元数据，发布时无需也无法指定。
- Computed 字段不破坏向后兼容，且不会因用户未设置而触发重建。

**替代方案：** 设为 `Optional + Computed` 并在 Create 传入 → 不采纳，云 API 创建接口不支持通过该参数指定来源版本（来源版本在 `CreateConfigGroupVersion` 时指定，而非发布时）。

### Decision 2: Read 回填沿用现有遍历与 nil 检查
**决策：** 在 Read 函数已有的 `for _, configGroupVersionInfo := range deployRecord.ConfigGroupVersionInfos` 循环中，新增一段：
```go
if configGroupVersionInfo.SourceVersion != nil {
    configGroupVersionInfoMap["source_version"] = configGroupVersionInfo.SourceVersion
}
```

**理由：**
- 与现有各字段（如 `VersionNumber`、`GroupId`）的回填方式完全一致，最小化改动。
- nil 检查遵循项目规范："在调用 setXX() 设置字段前，请先判断 Response 中的字段是否为 nil，若为 nil，则不调用 setXX()"。

### Decision 3: 不改动服务层
**决策：** `TeoService.DescribeTeoDeployConfigVersionHistoryByFilter` 返回 `[]*teov20220901.DeployRecord`，`DeployRecord.ConfigGroupVersionInfos` 已包含 `SourceVersion`，无需修改服务层。

**理由：** 服务层已透传整个 `ConfigGroupVersionInfo` 结构，新增字段自动可用，无需任何 SDK 升级。

## Risks / Trade-offs

### Risk 1: 云 API 未返回 SourceVersion 时字段为空
**风险：** 部分历史版本可能无来源版本，`SourceVersion` 为 nil。
**缓解措施：** Read 中保留 nil 检查，为 nil 时跳过 `set`，state 中该字段为空字符串，不影响资源 id 或其他字段。

### Risk 2: 旧 state 升级
**风险：** 已有资源 state 中不含 `source_version`，升级 provider 后首次 refresh 才填充。
**缓解措施：** Computed 字段缺失不影响 plan 结果，下次 `terraform refresh`/`apply` 自动回填，无破坏性变更。
