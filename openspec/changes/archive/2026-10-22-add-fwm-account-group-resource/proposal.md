## Why

云防火墙（FWM，Edge 边界防火墙账号组管理）目前没有 Terraform 资源用于管理「账号组」这一实体。账号组用于把集团成员账号归类，作为策略下发时的接收方集合。用户无法通过 Terraform 增删改查账号组，也无法以 IaC 方式描述账号组与集团策略之间的引用关系。新增 `tencentcloud_fwm_account_group` 资源填补该空白，使账号组纳入 Terraform 生命周期管理。

## What Changes

- 新增 Terraform 资源 `tencentcloud_fwm_account_group`，类型为 RESOURCE_KIND_GENERAL，覆盖完整 CRUD：
  - Create：调用 `CreateAccountGroup`，出参 GroupId 作为资源标识。
  - Read：服务端无单实例查询接口，调用 `DescribeAccountGroups` 并以 `Filters=[{Name:"GroupId",Values:[groupId],OperatorType:1}]` 精确取 `AccountGroups` 首元素；查不到（空数组）视为已删除，`d.SetId("")`。
  - Update：`group_name` 变更调用 `ModifyAccountGroup` 原地修改；`description` 变更因 `ModifyAccountGroup` 不接受该参数、仅 Create 可设置，标 `ForceNew` 触发重建。
  - Delete：调用 `DeleteAccountGroup`（软删，置 status=1）。
- 新增 FWM 服务目录 `tencentcloud/services/fwm/`，包含资源文件 `resource_tc_fwm_account_group.go`、服务层 `service_tencentcloud_fwm.go`、测试文件 `resource_tc_fwm_account_group_test.go`（使用 gomonkey mock，不使用 TF 测试套件）。
- 在 `tencentcloud/connectivity/client.go` 中新增 FWM 客户端 `UseFwmV20250611Client()` 及对应连接字段。
- 在 `tencentcloud/provider.go` 中注册 `tencentcloud_fwm_account_group`。
- 资源标识为服务端生成的 GroupId（`acg-xxxxxxxx`，8 位随机串），不是复合 ID；import 直接写 GroupId。
- 所有接口同步（无任务 ID、无中间态），不生成 waitReady/async 轮询。

## Capabilities

### New Capabilities
- `fwm-account-group`: 管理 FWM 账号组资源的完整生命周期（创建、读取、修改、删除、导入），对接 CreateAccountGroup/DescribeAccountGroups/ModifyAccountGroup/DeleteAccountGroup 四个云 API。

### Modified Capabilities

## Impact

- 新增文件：`tencentcloud/services/fwm/resource_tc_fwm_account_group.go`、`tencentcloud/services/fwm/service_tencentcloud_fwm.go`、`tencentcloud/services/fwm/resource_tc_fwm_account_group_test.go`、`tencentcloud/services/fwm/resource_tc_fwm_account_group.md`。
- 修改文件：`tencentcloud/provider.go`（注册资源 + import）、`tencentcloud/provider.md`（注册资源）、`tencentcloud/connectivity/client.go`（新增 UseFwmV202506Client 及连接字段）。
- 云 API 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（已在 vendor/modules.txt 中存在）。
- 与 `tencentcloud_cfw_*` 资源分属不同的产品线（cfw 为 v20190904，fwm 为 v20250611），各自独立的服务目录与客户端，互不影响。