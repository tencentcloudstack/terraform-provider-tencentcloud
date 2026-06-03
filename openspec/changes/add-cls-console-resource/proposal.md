## Why

CLS DataSight 控制台目前没有 Terraform 管理通道，用户只能在网页控制台创建/修改/删除 DataSight 实例，无法做基础设施代码化（IaC）。需要新增一个 CRUD 型 Terraform 资源，把 CLS DataSight 控制台 4 个新接口（CreateConsole / DescribeConsoles / ModifyConsole / DeleteConsole）暴露为 `tencentcloud_cls_console`，让用户能在 HCL 中声明 DataSight 实例并通过 `terraform plan/apply` 管理其全生命周期。

## What Changes

- 新增资源 `tencentcloud_cls_console`（CRUD 型）。
- 新增 Schema 字段，**严格按 CreateConsole 接口入参 1:1 映射**：`access_mode`、`login_mode`、`domain_prefix`、`accounts`、`anonymous_login`、`intranet_type`、`intranet_region`、`vpc_id`、`subnet_id`、`auth_roles`、`tags`、`hide_params`、`access_control_rules`、`remarks`、`menus`。
- 新增 Computed 字段：`console_id`、`domain`、`intranet_domain`（来自 Describe 响应，由后端生成）。
- 资源 ID 设为 `ConsoleId`（`d.SetId(consoleId)`）。
- 新增 service 层方法：`DescribeClsConsoleById(ctx, consoleId)`，内部调用 `DescribeConsoles` + `Filters` 过滤 `ConsoleId`，分页 limit 取接口最大值 100。
- Update 路径调用 `ModifyConsole`（注意：ModifyConsole 不接受 `Tags`，所以 `tags` 字段为 ForceNew 或在 Update 中被忽略——本 change 选择**ForceNew**，与 Tencent 其他 cls 资源一致）。
- Delete 路径调用 `DeleteConsole`，请求体仅 `ConsoleId`。
- 所有 SDK 调用必须用 `resource.Retry(...)` 包装。
- 新增 `.md` 资源文档（命名与 `resource_tc_config_compliance_pack.md` 同款）+ 验收测试 `_test.go`（命名与 `resource_tc_config_compliance_pack_test.go` 同款）。
- 在 `tencentcloud/provider.go` 注册新资源 `tencentcloud_cls_console`。

## Capabilities

### New Capabilities

- `cls-console-resource`: 新增 `tencentcloud_cls_console` 资源的 schema、CRUD 行为、ID 约定、字段约束、文档与测试规范，作为 CLS DataSight 控制台的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go 增加注册一行 -->

## Impact

- 新文件：
  - `tencentcloud/services/cls/resource_tc_cls_console.go`
  - `tencentcloud/services/cls/resource_tc_cls_console.md`
  - `tencentcloud/services/cls/resource_tc_cls_console_test.go`
  - `website/docs/r/cls_console.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/services/cls/service_tencentcloud_cls.go`：新增 `DescribeClsConsoleById` 方法（不破坏现有方法）
  - `tencentcloud/provider.go`：在资源注册映射中追加 `"tencentcloud_cls_console": cls.ResourceTencentCloudClsConsole()` 一行
- 依赖：现有 SDK `tencentcloud-sdk-go/tencentcloud/cls/v20201016` 已包含 4 个所需接口（`CreateConsole/DescribeConsoles/ModifyConsole/DeleteConsole`）以及配套 model（`Console`、`ConsoleAccount`、`AnonymousLoginInfo`、`AuthRoleInfo`、`AccessControlRule`、`Tag`、`Filter`），无需升级 SDK。
- 不修改任何既有资源的 schema，对现有用户零影响。
