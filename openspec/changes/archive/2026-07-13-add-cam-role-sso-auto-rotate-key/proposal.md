## Why

`tencentcloud_cam_role_sso` 资源当前仅支持配置 OIDC 身份提供商的 `name`、`identity_url`、`identity_key`、`client_ids`、`description` 五个参数，缺少 OIDC 公钥自动轮转开关的管理能力。腾讯云 CAM 的 OIDC 配置接口（CreateOIDCConfig / UpdateOIDCConfig / DescribeOIDCConfig）已支持 `AutoRotateKey` 参数，用户无法通过 Terraform 自动化开启公钥自动轮转，限制了安全运维的自动化水平。本次变更新增 `auto_rotate_key` 参数，让用户能够以声明式方式管理 OIDC 公钥的自动轮转开关。

## What Changes

- 在 `tencentcloud_cam_role_sso` 资源 schema 中新增可选参数 `auto_rotate_key`（类型 `schema.TypeInt`，OIDC 公钥自动轮转开关，枚举值 0=关闭、1=开启，默认值 0）。
- 在 `resourceTencentCloudCamRoleSSOCreate` 中将 `auto_rotate_key` 透传到 `CreateOIDCConfigRequest.AutoRotateKey`。
- 在 `resourceTencentCloudCamRoleSSOUpdate` 中将 `auto_rotate_key` 纳入变更检测并透传到 `UpdateOIDCConfigRequest.AutoRotateKey`。
- 在 `resourceTencentCloudCamRoleSSORead` 中从 `DescribeOIDCConfigResponse.AutoRotateKey` 回填 `auto_rotate_key` 到 state（回填前判断字段是否为 nil）。
- 在 `resource_tc_cam_role_sso.md` 文档的示例与 import 说明中补充新参数。

## Capabilities

### New Capabilities
- `cam-role-sso-resource`: 新增能力，描述 `tencentcloud_cam_role_sso` 资源对 OIDC 公钥自动轮转开关（`auto_rotate_key`）参数的管理行为（创建、更新、读取、删除时的处理）。

### Modified Capabilities
<!-- 无既有 spec 需要修改。openspec/specs/ 中不存在 cam-role-sso 相关 spec，因此本次以新增能力 spec 的形式落地。 -->

## Impact

### 受影响的代码
- `tencentcloud/services/cam/resource_tc_cam_role_sso.go` - 新增 `auto_rotate_key` schema 字段，并在 Create/Read/Update 函数中处理该参数。
- `tencentcloud/services/cam/resource_tc_cam_role_sso_test.go` - 补充针对 `auto_rotate_key` 的单元测试用例（使用 gomonkey mock 云 API）。
- `tencentcloud/services/cam/resource_tc_cam_role_sso.md` - 更新资源文档示例。

### 云 API 依赖
基于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116`（vendor 中已存在）：
- `CreateOIDCConfigRequest.AutoRotateKey`（`*uint64`）- 创建时入参 ✓
- `UpdateOIDCConfigRequest.AutoRotateKey`（`*uint64`）- 更新时入参 ✓
- `DescribeOIDCConfigResponse.AutoRotateKey`（`*uint64`）- 查询时出参 ✓
- `DeleteOIDCConfigRequest` 仅有 `Name` 字段，删除操作不涉及 `AutoRotateKey` ✓

### 向后兼容性
- ✅ 完全向后兼容，仅新增一个 Optional 字段，不修改已有 schema 行为。
- ✅ 现有不带 `auto_rotate_key` 的 TF 配置继续正常工作，API 默认值（0）生效。
- ✅ 无 state 迁移，无破坏性变更。

### 测试影响
- 新增参数需要补充单元测试（mock 云 API），覆盖创建、更新、读取场景下 `auto_rotate_key` 的设置与回填。
