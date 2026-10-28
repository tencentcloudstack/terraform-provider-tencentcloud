## Context

`tencentcloud_clb_target_group` 资源在过往迭代中已经将 CLB 目标组的核心字段（health_check、schedule_algorithm、weight、session_expire_time、ip_version 等）打通到 Terraform schema。SDK 侧 (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`) 的三个关键 API 均已包含 `SnatEnable *bool`：

- `CreateTargetGroupRequest.SnatEnable`（`models.go` ≈1971 行）
- `ModifyTargetGroupAttributeRequest.SnatEnable`（`models.go` ≈8094 行）
- `TargetGroupInfo.SnatEnable`（响应结构，`DescribeTargetGroupList` 返回的 `TargetGroupSet` 元素）

但 Provider 层并未透传该字段，导致用户无法通过 IaC 管理 SNAT 行为。

约束：
- 严格保持向后兼容：仅新增 Optional 字段，不修改已有字段语义。
- 不引入新的 SDK 依赖或 API。
- 同步更新 `.md` 资源样例并通过 `make doc` 重新生成 `website/docs/`。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_clb_target_group` 资源中提供 `snat_enable` 参数，覆盖 Create / Read / Update 三个生命周期阶段。
- service 层 `ClbService.CreateTargetGroup` 在不破坏既有调用方语义的前提下接收新的 `snatEnable *bool` 入参。
- 验收测试覆盖：创建时启用 SNAT、创建时关闭 SNAT、Update 切换 SNAT 状态。
- 文档与示例同步更新。

**Non-Goals:**
- 不修改 `tencentcloud_clb_target_group_attachment`、`tencentcloud_clb_target_group_instance_attachment` 等关联资源（这些资源不直接持有 SnatEnable）。
- 不更改 `data.tencentcloud_clb_target_groups` / `data.tencentcloud_clb_target_group_list` 的输出 schema（如需后续暴露，可在独立变更中单独评估，本变更聚焦资源 CRUD）。
- 不更改 SDK 版本。
- 不为 `snat_enable` 设置 `Default`：保持腾讯云 API 默认行为（关闭），由云端决定。

## Decisions

### Decision 1：Schema 字段使用 `Optional + Computed`，不设 `Default`
- **选择**：`Type: schema.TypeBool, Optional: true, Computed: true`，不设置 `Default`。
- **理由**：
  - 与同资源已有的 `keepalive_enable`、`full_listen_switch` 等布尔开关风格一致。
  - 若设置 `Default: false`，会导致存量用户在 `terraform plan` 时被动出现 diff（state 中没有该字段，计划将其显式置为 false 并下发 API）。
  - `Computed` 允许从 API 回填真实值，避免与云端默认值不一致带来的漂移。
- **替代方案**：仅 `Optional`、不 `Computed` —— 缺点是初次 import 或资源由控制台先创建后导入时无法回写真实状态。
- **替代方案**：`Optional + Default: false` —— 缺点见上，会引入兼容性问题。

### Decision 2：`ForceNew` 与否
- **选择**：不设置 `ForceNew`，因为 `ModifyTargetGroupAttribute` SDK 支持在线修改 `SnatEnable`。
- **理由**：与 SDK 能力对齐，避免不必要的资源重建。
- **风险**：腾讯云对 SnatEnable 在某些目标组类型（v1 vs v2）下的修改限制可能存在差异。
  - **缓解**：在 Update 函数中将 API 返回的错误如实抛出，由用户根据云端报错调整；同时在文档中加入注意事项。

### Decision 3：service 层 `CreateTargetGroup` 函数签名扩展
- **选择**：在末尾追加 `snatEnable *bool` 参数（在 `ipVersion string` 之后）。
- **理由**：
  - 现有所有调用点都在本仓库内（`resource_tc_clb_target_group.go` 唯一调用），可统一更新。
  - 使用 `*bool` 而非 `bool`，与已有 `fullListenSwitch *bool` / `keepaliveEnable *bool` 风格一致，nil 表示不传该字段，区分 "未设置" 与 "显式 false"。
- **替代方案**：构造一个 options struct 重构所有可选项 —— 范围太大，超出本变更目标，留待后续独立的 refactor。
- **替代方案**：在 service 内部从 `*schema.ResourceData` 取值 —— 违反职责分层，不采用。

### Decision 4：Update 流程中的 diff 检测
- **选择**：使用 `d.HasChange("snat_enable")` + `d.GetOkExists("snat_enable")` / `d.Get` 取值，并将结果赋给 `request.SnatEnable`。
- **理由**：
  - `HasChange` 能正确感知从 nil → false / true 等所有过渡。
  - 直接使用 `d.Get("snat_enable").(bool)` 即可（`Computed` 字段不会出现 `nil`），通过 `helper.Bool()` 包装。
- **细节**：在 Update 中，将 SNAT 改动加入 `isChanged = true` 触发器分支，与 `keepalive_enable` 等同类字段并列。

### Decision 5：Read 流程回写
- **选择**：在 `resourceTencentCloudClbTargetRead` 中，若 `targetGroup.SnatEnable != nil` 则 `_ = d.Set("snat_enable", *targetGroup.SnatEnable)`，与同函数中其他 nil 检查模式保持一致。

### Decision 6：测试策略
- **选择**：新增独立验收测试 `TestAccTencentCloudClbTargetGroup_snatEnable`，使用两步 `TestStep`：
  - Step 1：创建时 `snat_enable = true`，校验 state 与远端一致。
  - Step 2：Update 为 `snat_enable = false`，校验更新成功且未触发资源重建。
- 同时在 `.md` 示例文件中给出展示性配置。
- 受限于验收测试需要真实账号，PR 提交时附运行截图/日志即可。

## Risks / Trade-offs

- **风险**：腾讯云对 `SnatEnable` 与 v1/v2 目标组、协议组合可能存在隐藏限制（例如某些 v1 目标组不支持 SNAT）。
  → **缓解**：依赖云 API 自身校验，错误信息透传给用户；在文档中说明 "SnatEnable 是否生效取决于目标组类型与协议，由云端决策"。
- **风险**：`Computed` 字段在 API 偶发返回不一致时可能导致漂移。
  → **缓解**：与 `keepalive_enable` 等已上线字段处理方式一致，依赖 `defer tccommon.InconsistentCheck()` 兜底。
- **取舍**：未在数据源 `tencentcloud_clb_target_groups` 中输出 `snat_enable`。
  → **理由**：避免本变更范围扩散；后续如有需求可单独评估。

## Migration Plan

- 此变更纯增量，无需迁移。
- 存量 `tencentcloud_clb_target_group` 资源在升级后首次 `terraform refresh` 会通过 `Computed` 路径自动回填 `snat_enable`，不会产生 diff。
- 回滚：直接还原 schema 与 service 层改动；已写入 state 的 `snat_enable` 字段在旧版本下会被忽略，不影响功能。

## Open Questions

- 是否需要在数据源（`tencentcloud_clb_target_groups` 与 `tencentcloud_clb_target_group_list`）中也暴露 `snat_enable`？  
  → 暂列为非目标，待用户反馈再决定，独立变更跟进。
