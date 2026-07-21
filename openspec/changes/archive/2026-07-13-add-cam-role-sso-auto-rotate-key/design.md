## Context

`tencentcloud_cam_role_sso` 是一个 RESOURCE_KIND_GENERAL 资源，管理腾讯云 CAM 的角色 SSO（OIDC）配置，位于 `tencentcloud/services/cam/resource_tc_cam_role_sso.go`。其 CRUD 分别对应：
- Create → `CreateOIDCConfig`
- Read → `DescribeOIDCConfig`
- Update → `UpdateOIDCConfig`
- Delete → `DeleteOIDCConfig`

当前 schema 仅包含 `name`、`identity_url`、`identity_key`、`client_ids`、`description` 五个字段。腾讯云 CAM SDK（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116`）中 `CreateOIDCConfigRequest`、`UpdateOIDCConfigRequest`、`DescribeOIDCConfigResponse` 均已新增 `AutoRotateKey`（`*uint64`）字段，表示 OIDC 公钥自动轮转开关（0=关闭，1=开启，默认 0）。本次变更将这个云 API 能力透传到 Terraform 资源。

约束：
- 必须保持向后兼容，只能新增 Optional 字段。
- 遵循项目代码规范：retry 块内只调用云 API，错误用 `tccommon.RetryError()` 包装；Read 中 set 前判断 nil；Create 后检查返回值。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_cam_role_sso` 资源中新增 `auto_rotate_key` 参数，支持创建、更新、读取该字段。
- 保证现有配置与 state 不受影响（向后兼容）。
- 补充单元测试（gomonkey mock 云 API），覆盖新参数的创建/更新/读取逻辑。

**Non-Goals:**
- 不修改 `DeleteOIDCConfig` 调用逻辑（`DeleteOIDCConfigRequest` 不支持 `AutoRotateKey`）。
- 不新增/修改其他 CAM 资源或数据源。
- 不重构现有 schema 字段。
- 不在本阶段执行 `gofmt`、`make doc`、`.changelog` 文件创建（统一由收尾阶段 tfpacer-finalize skill 处理）。

## Decisions

### Decision 1: 字段类型选择 `schema.TypeInt`
云 API 中 `AutoRotateKey` 为 `*uint64`，取值枚举为 0/1。选择 `schema.TypeInt` 而非 `schema.TypeBool`：
- 与同 SDK 中 `CreateUserOIDCConfigRequest.AutoRotateKey`（同样是 `*uint64`，0/1 枚举）保持一致，避免布尔与数值的转换歧义。
- 直接与云 API 的数值语义对应，`helper.IntUint64()` / `helper.IntUint64Ptr()` 进行 int↔uint64 转换。
- **替代方案**：使用 `schema.TypeBool`（true=1, false=0）。放弃，因为云 API 明确使用 0/1 枚举且默认值文档化，TypeInt 语义更直观，且资源文件中已有其他数值开关字段的惯例。

### Decision 2: 字段为 Optional、非 ForceNew
- `auto_rotate_key` 是可变更的运行时开关，应在 Update 中支持修改，因此不设 ForceNew。
- 在 `resourceTencentCloudCamRoleSSOUpdate` 的 `d.HasChange(...)` 检测条件中加入 `auto_rotate_key`，使其变更触发 `UpdateOIDCConfig` 调用。

### Decision 3: Create 中始终透传 auto_rotate_key
在 `resourceTencentCloudCamRoleSSOCreate` 中，无论用户是否配置 `auto_rotate_key`，都将其从 schema 读取并设置到 `request.AutoRotateKey`。未配置时 schema 默认返回 0，与云 API 默认值一致，行为安全。这样实现简洁，且不会因为 omit 导致与用户预期不一致。

### Decision 4: Read 中 nil 判断后回填
`DescribeOIDCConfigResponse.Response.AutoRotateKey` 可能为 nil。遵循项目规范，在 `resourceTencentCloudCamRoleSSORead` 中先判断 `response.Response.AutoRotateKey != nil` 再 `d.Set("auto_rotate_key", *response.Response.AutoRotateKey)`，避免空指针。

### Decision 5: 单元测试使用 gomonkey mock
本资源为现有资源修改，按项目要求，现有资源的测试用例补充使用 Terraform 测试套件。但本次同时需用 `go test -gcflags=all=-l` 跑通涉及文件——为兼顾"新增参数"的可独立验证性，在 `resource_tc_cam_role_sso_test.go` 中使用 gomonkey 对 `UseCamClient().CreateOIDCConfig` / `UpdateOIDCConfig` / `DescribeOIDCConfig` / `DeleteOIDCConfig` 进行 mock，编写纯业务逻辑单元测试（不依赖 TF_ACC），验证 `auto_rotate_key` 在 Create/Update/Read 中的正确传递与回填。

## Risks / Trade-offs

- [风险] 云 API `DescribeOIDCConfig` 旧版本可能不返回 `AutoRotateKey`（返回 nil） → Read 中已加 nil 判断，跳过 set，state 中保留 schema 默认值，不影响现有资源。
- [风险] 用户在配置中显式写 `auto_rotate_key = 0` 与不写的语义差异 → 两者行为一致（均传 0 给 API），无副作用。
- [权衡] 使用 TypeInt 而非 TypeBool → 与云 API 数值语义一致，但用户需了解 0/1 枚举含义，已在 Description 中说明。

## Migration Plan

无需迁移。新增 Optional 字段，现有配置不受影响，直接部署即可。回滚仅需移除该 schema 字段及相关 CRUD 代码。
