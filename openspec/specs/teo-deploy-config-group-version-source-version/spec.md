# teo-deploy-config-group-version-source-version Specification

## Purpose
Defines that the `tencentcloud_teo_deploy_config_group_version` resource SHALL expose a `source_version` field within its `config_group_version_infos` nested block, mapping to the TencentCloud EdgeOne `ConfigGroupVersionInfo.SourceVersion` field returned by the `DescribeDeployHistory` API. The field is read-only and surfaces the source version ID from which a version was derived.

## Requirements
### Requirement: teo_deploy_config_group_version 暴露 source_version 字段
`tencentcloud_teo_deploy_config_group_version` 资源的 `config_group_version_infos` 嵌套块 SHALL 新增 `source_version` 字段，用于映射 EdgeOne 云 API `ConfigGroupVersionInfo.SourceVersion`（版本派生的来源版本 ID）。

#### Scenario: Schema 定义 source_version 为 Computed
- **GIVEN** `tencentcloud_teo_deploy_config_group_version` 资源的 `config_group_version_infos` 嵌套 schema
- **WHEN** 检查 `source_version` 字段属性
- **THEN** 该字段 SHALL 为 `Type: schema.TypeString` 且 `Computed: true`
- **AND** 该字段 SHALL 不设 `Required`、`Optional` 或 `ForceNew`

#### Scenario: Read 回填 source_version
- **GIVEN** 一个已创建的 `tencentcloud_teo_deploy_config_group_version` 资源
- **WHEN** Read 函数调用 `DescribeDeployHistory` 且响应 `Records[].ConfigGroupVersionInfos[].SourceVersion` 非空
- **THEN** Provider SHALL 将 `SourceVersion` 值写入 state 中 `config_group_version_infos[].source_version`
- **AND** 回填前 SHALL 判断 `SourceVersion` 是否为 nil，为 nil 时 SHALL 不执行写入

#### Scenario: 云 API 未返回 SourceVersion
- **GIVEN** `DescribeDeployHistory` 响应中某条 `ConfigGroupVersionInfos` 的 `SourceVersion` 为 nil
- **WHEN** Read 函数处理该条记录
- **THEN** Provider SHALL 跳过 `source_version` 字段的写入
- **AND** SHALL 不返回错误， SHALL 不清空资源 id

### Requirement: Create 流程不使用 source_version
`source_version` SHALL 仅作为只读展示字段，资源 Create 流程 SHALL NOT 将其作为 `DeployConfigGroupVersion` 请求的入参。

#### Scenario: Create 不传递 source_version
- **GIVEN** 用户在 Terraform 配置中（Computed 字段不可配置，故无法显式设置）
- **WHEN** 执行 `terraform apply` 创建资源
- **THEN** `DeployConfigGroupVersion` 请求的 `ConfigGroupVersionInfos[].SourceVersion` SHALL NOT 由该字段设置
- **AND** Create 请求 SHALL 仅设置 `VersionId`（保持现有行为）

### Requirement: 文档更新
资源文档 SHALL 描述新增的 `source_version` 字段。

#### Scenario: 资源 md 文档包含 source_version
- **GIVEN** `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.md`
- **WHEN** 查看 `config_group_version_infos` 嵌套块的字段说明
- **THEN** 文档 SHALL 包含 `source_version` 字段
- **AND** 说明其为来源版本 ID（只读）
