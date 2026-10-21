## Why

腾讯云防火墙（FWM，API 版本 `2025-06-11`）的「规则组」目前没有 Terraform 管理通道：用户只能在控制台手工创建/改名/删除规则组，无法把安全策略的骨架（规则组）纳入 IaC 流程。`fwm` 的 `CreateRuleGroup` / `DescribeRuleGroups` / `ModifyRuleGroup` / `DeleteRuleGroup` 四个接口均为**同步接口**，语义完整、无异步任务 ID，适合直接映射为一个 Terraform CRUD 资源。

## What Changes

- 新增 RESOURCE_KIND_GENERAL 资源 `tencentcloud_fwm_rule_group`，归属新建服务包 `tencentcloud/services/fwm/`，实现完整 CRUD（资源生命周期全覆盖）。
- 文件命名遵循项目规范：`resource_tc_fwm_rule_group.go` / `resource_tc_fwm_rule_group.md` / `resource_tc_fwm_rule_group_test.go`。
- Schema 严格来自云 API 可写字段：`group_name`（Create/Modify 均支持）、`product`（仅 Create 支持，**ForceNew**）；其余 `RuleGroupResp` 返回字段按项目规范**平铺**为 Computed 字段（`app_id`、`rule_count`、`create_time`、`update_time`、`create_by`、`create_name`、`update_by`、`update_name`、`receive_account`、`exec_areas`、`direction`）。**不创建** `rule_groups` 这类「列表型数据」包装层。
- 资源 ID = `CreateRuleGroup` 返回的 `GroupId`（云 API 唯一标识，`fwmrg-` 前缀），单字段 ID，不需要复合 ID。
- Read 通过 `DescribeRuleGroups` + `Filters=[{Name:"GroupId", Values:[<id>], OperatorType:9}]` 精确定位单条记录（不使用 `Product` 过滤，避免 `Product` 入/出参类型不一致导致查不到）。
- **`Product` 类型归一化**：`CreateRuleGroup` 入参为字符串（`enterprise_sg` / `cfw_nat_acl` / `cfw_vpc_acl` / `cfw_edge_acl`），而 `DescribeRuleGroups` 云 API 实际返回数字（`1` / `2` / `3` / `4`）。Read 时把数字映射回字符串写入 state，避免 `terraform plan` 永久 diff。
- Update 仅支持 `group_name`（`ModifyRuleGroup` 只有 `GroupId` + `GroupName`）；`Product` 在创建时固定，标记 `ForceNew`。
- Delete 调用 `DeleteRuleGroup`，`GroupIds` 传单个元素 `[<groupId>]`。
- 新增 service 层 `DescribeFwmRuleGroupById`；新增 connectivity 客户端工厂 `UseFwmV20250611Client()`。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 注册新资源。
- 新增基于 gomonkey 的单元测试（新资源不写 acceptance test 套件），覆盖 Create/Read/Update/Delete 及 `Product` 数字归一化。

## Capabilities

### New Capabilities

- `fwm-rule-group-resource`: 定义 `tencentcloud_fwm_rule_group` 的 schema 字段集、资源 ID 约定、CRUD 行为（同步接口、无轮询）、`Product` 枚举与数字归一化规则、文档与单元测试要求。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go / provider.md 增加注册项 -->

## Impact

- 新增文件：
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_rule_group_test.go`
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`
  - `website/docs/r/fwm_rule_group.html.markdown`（收尾阶段由 `make doc` 生成，禁止手写）
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `fwmV20250611Conn` 字段与 `UseFwmV20250611Client()`（纯新增，不影响既有服务）
  - `tencentcloud/provider.go`：新增 import 与资源注册一行
  - `tencentcloud/provider.md`：末尾新增 `Firewall Manager(FWM)` 段与资源名
- 依赖：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（已 vendor，`vendor/modules.txt` 已包含），无需升级 SDK。
- 不修改任何既有资源的 schema，对现有用户零影响。
