# teo-deploy-config-group-version-env-info Specification

## Purpose
TBD - created by archiving change add-teo-deploy-config-group-version-env-info-params. Update Purpose after archive.
## Requirements
### Requirement: Read environment-level info via DescribeEnvironments

`resourceTencentCloudTeoDeployConfigGroupVersionRead` 在解析出 `zone_id` 与 `env_id` 后，除读取部署记录（`DescribeConfigGroupVersionHistory`）外，SHALL 额外调用 `DescribeEnvironments`（复用现有 service 层方法 `DescribeTeoEnvironmentsByFilter`，入参传入 `ZoneId`）定位 `EnvId == env_id` 的 `EnvInfo`，并将其字段映射为资源的 Computed 出参写入 state，以便用户跟踪版本发布后环境的实际生效状态与当前生效配置组版本信息。

#### Scenario: 环境存在时填充全部新增出参

- **WHEN** `DescribeEnvironments` 返回的 `EnvInfos` 中存在 `EnvId` 等于资源 `env_id` 的元素，且该元素各字段非空
- **THEN** 系统 SHALL 将以下字段写入对应 Computed schema：
  - `EnvInfo.EnvType` → `env_type`
  - `EnvInfo.Scope`（`[]*string`）→ `scope`（string 列表，按原顺序保留元素）
  - `EnvInfo.CurrentConfigGroupVersionInfos`（`[]*ConfigGroupVersionInfo`）→ `current_config_group_version_infos`（集合），其中每个元素映射：
    - `VersionId` → `version_id`
    - `VersionNumber` → `version_number`
    - `SourceVersion` → `source_version`
    - `GroupType` → `group_type`
    - `GroupId` → `group_id`
    - `Description` → `description`
    - `Status` → `status`
    - `CreateTime` → `create_time`
  - `EnvInfo.CreateTime` → `env_create_time`
  - `EnvInfo.UpdateTime` → `env_update_time`

#### Scenario: DescribeEnvironments 返回环境总数

- **WHEN** `DescribeEnvironments` 响应 `Response.TotalCount` 非空
- **THEN** 系统 SHALL 将 `TotalCount`（`uint64`）写入资源 Computed 出参 `total_count`

#### Scenario: 目标环境不存在时不报错且不影响既有部署记录读取

- **WHEN** `DescribeEnvironments` 返回的 `EnvInfos` 中没有任何元素的 `EnvId` 等于资源 `env_id`（例如返回为空列表或仅含其它环境）
- **THEN** 系统 SHALL 不对该新增出参做任何赋值（保持未设置），且 SHALL 不影响既有部署记录（`record_id`、`deploy_time`、`status`、`message`、`config_group_version_infos`）的读取，整个 Read 方法 SHALL 正常返回（不返回错误）

#### Scenario: 字段为空时跳过对应 set 调用

- **WHEN** 定位到的 `EnvInfo` 中某个字段为 `nil`（例如 `Scope` 为 nil、`CurrentConfigGroupVersionInfos` 为 nil、`CreateTime` 为 nil）
- **THEN** 系统 SHALL 跳过该字段的 `d.Set()` 调用，不对 nil 字段强行写入，其余非空字段仍按 Scenario "环境存在时填充全部新增出参" 写入

### Requirement: New computed schema fields are read-only and backward compatible

`tencentcloud_teo_deploy_config_group_version` 资源 SHALL 新增以下 Computed 字段：`total_count`、`env_type`、`scope`、`current_config_group_version_infos`（嵌套集合，含 `version_id` / `version_number` / `source_version` / `group_type` / `group_id` / `description` / `status` / `create_time`）、`env_create_time`、`env_update_time`。所有新增字段 MUST 为 `Computed: true` 且不可由用户输入（不设 `Required`/`Optional`），从而保持向后兼容：现有 Terraform 配置与 state 不受影响，`terraform plan` 对既有配置不产生 diff。

#### Scenario: 新增字段为 Computed 且不可配置

- **WHEN** 用户编写 `tencentcloud_teo_deploy_config_group_version` 资源的 HCL 配置
- **THEN** `total_count`、`env_type`、`scope`、`current_config_group_version_infos`、`env_create_time`、`env_update_time` SHALL 不出现在用户可配置参数中（均为 Computed-only），用户无需也无法显式设置这些字段

#### Scenario: 向后兼容既有配置

- **WHEN** 一个在本次变更前创建的 `tencentcloud_teo_deploy_config_group_version` 资源配置在本次变更后再次执行 `terraform plan`
- **THEN** 该配置 SHALL 不产生由本次新增字段导致的 diff（新增字段均为 Computed，不影响既有 `zone_id` / `env_id` / `config_group_version_infos` / `description` 等字段语义）

