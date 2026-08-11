## Why

`tencentcloud_teo_config_group_version` 资源用于管理 EdgeOne (TEO) 配置组版本，资源复合 ID 由 `ZoneId`、`GroupId`、`VersionId` 拼接而成。当前 `DescribeConfigGroupVersionDetail` 云 API 的 `Response.ConfigGroupVersionInfo.VersionId` 出参未被规范地暴露为独立的 Terraform schema 参数，导致用户无法在 state 中明确读取到配置组版本的 VersionId 标识，不利于跨资源引用与排障。腾讯云 TEO SDK 已提供该出参，需要将其补充为资源的顶层输出参数。

## What Changes

- 在 `tencentcloud_teo_config_group_version` 资源中新增出参 `version_id`（Computed, TypeString），来源于 `DescribeConfigGroupVersionDetail` 接口返回的 `Response.ConfigGroupVersionInfo.VersionId`。
- 在资源 Read 方法中，从 `respData.ConfigGroupVersionInfo.VersionId` 读取并 `d.Set("version_id", ...)` 回填到 Terraform state（仅在字段非 nil 时设置）。
- 更新资源文档 `resource_tc_teo_config_group_version.md`，补充对 `version_id` 出参的说明。
- 在单元测试 `resource_tc_teo_config_group_version_test.go` 中补充覆盖 `version_id` 出参读取的用例（使用 gomonkey mock 云 API）。

非破坏性变更：`version_id` 为新增的 Computed 出参，不影响现有 TF 配置与 state。

## Capabilities

### New Capabilities
- `teo-config-group-version-resource`: 管理 EdgeOne (TEO) 配置组版本资源的 schema 定义与 Read 行为，包含从 `DescribeConfigGroupVersionDetail` 读取 `version_id` 出参的能力。

### Modified Capabilities
<!-- 无既有 capability 被修改 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version.go`（schema 新增 `version_id` 出参声明，Read 方法设置该字段）
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version_test.go`（补充 `version_id` 出参读取的单测用例）
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version.md`（补充 `version_id` 出参说明）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.DescribeConfigGroupVersionDetailResponse`，无需变更 vendor。
- 向后兼容：新增 Computed 出参，不影响已有 state 与 TF 配置（Computed 字段由 Read 回填，不会触发 plan diff）。
- 文档：需要同步更新 `website/docs/` 下文档（由 `make doc` 自动生成流程读取 `.md` 文件生成，禁止手改 website 目录）。
