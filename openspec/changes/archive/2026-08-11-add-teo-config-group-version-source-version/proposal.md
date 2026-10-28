## Why

`tencentcloud_teo_config_group_version` 资源用于在 EdgeOne（TEO）中创建配置组版本。当前资源在创建版本时只能基于当前生产环境生效的版本派生新版本，用户无法指定来源版本 ID（`SourceVersion`）。腾讯云 TEO SDK 的 `CreateConfigGroupVersion` 接口已支持可选入参 `SourceVersion`，且 `DescribeConfigGroupVersionDetail` 接口的返回值 `ConfigGroupVersionInfo` 中也包含 `SourceVersion` 字段，但 Terraform Provider 未暴露该参数，导致用户无法通过声明式配置指定新版本所基于的来源版本，限制了版本派生管理的灵活性。

## What Changes

- 在 `tencentcloud_teo_config_group_version` 资源 Schema 中新增 `source_version` 字段：
  - 类型：`schema.TypeString`
  - 可选：`Optional: true`，`ForceNew: true`（资源仅支持 Create/Read/Delete，无 Update 接口）
  - 创建时传递给 `CreateConfigGroupVersion` API 的 `request.SourceVersion`
  - 读取时从 `DescribeConfigGroupVersionDetail` 返回的 `ConfigGroupVersionInfo.SourceVersion` 回填到状态
- 在资源 Create 方法中新增对 `source_version` 参数的读取与请求构造
- 在资源 Read 方法中新增对 `ConfigGroupVersionInfo.SourceVersion` 的 nil 检查与 `d.Set` 回填
- 在单元测试文件 `resource_tc_teo_config_group_version_test.go` 中补充覆盖 `source_version` 的测试用例
- 更新资源文档 `resource_tc_teo_config_group_version.md`，在示例中展示 `source_version` 字段

非破坏性：新增 Optional 字段，未配置时使用云平台默认行为（当前生产环境生效版本），不影响已有 TF 配置与 state。

## Capabilities

### New Capabilities
- `teo-config-group-version-resource`: TEO 配置组版本（ConfigGroupVersion）资源管理，涵盖 Create/Read/Delete 操作及 `source_version` 等参数的 Schema 定义与 API 映射

### Modified Capabilities
<!-- 无现有 capability 需修改 -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version.go`（Schema 新增字段、Create/Read 逻辑补充）
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version_test.go`（新增 source_version 测试用例）
  - `tencentcloud/services/teo/resource_tc_teo_config_group_version.md`（示例补充 source_version）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.CreateConfigGroupVersionRequest.SourceVersion` 与 `teov20220901.ConfigGroupVersionInfo.SourceVersion`，无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置；未设置 `source_version` 时保持原有行为（由云平台默认采用当前生产环境生效版本）。
- 文档：需要同步更新 website docs（由收尾阶段 `make doc` 自动生成流程读取 `.md` 文件）。
