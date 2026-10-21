## Why

云防火墙 FWM（FireWall Management）的 VPC 边界 ACL 规则目前没有 Terraform 管理通道。用户只能在网页控制台逐条创建/修改/删除 VPC 边界规则，无法做基础设施代码化（IaC）。需要新增一个 CRUD 型 Terraform 资源，把 FWM 的 4 个同步接口（CreateVpcAclRule / DescribeVpcAclRules / ModifyVpcAclRule / DeleteVpcAclRule）暴露为 `tencentcloud_fwm_vpc_acl_rule`，让用户能在 HCL 中声明单条 VPC 边界 ACL 规则并通过 `terraform plan/apply` 管理其全生命周期。

VPC 边界规则组是父容器（由独立资源 fwm_rule_group 创建），本资源只管理「单条 VPC 边界 ACL 规则」，资源 ID 为复合结构 `GroupId#RuleId`。

## What Changes

- 新增资源 `tencentcloud_fwm_vpc_acl_rule`（RESOURCE_KIND_GENERAL，CRUD 型）。
- 新增 service 层 `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（fwm 是全新服务目录，需同时创建 client 接入）。
- 新增 `tencentcloud/connectivity/client.go` 中 `UseFwmV20250611Client()` 访问器与 `fwmv20250611Conn` 字段、import。
- 新增 Schema 字段：
  - 必填字段：`group_id`（ForceNew，父资源引用）、`edge_id`（ForceNew）、`fw_group_id`（ForceNew）、`ip_version`（ForceNew）、`source_type`（ForceNew）、`source_content`（ForceNew）、`dest_type`（ForceNew）、`dest_content`（ForceNew）、`protocol`、`rule_action`、`order_index`、`description`（必填，≤30 字符）。
  - 可选字段：`port`、`param_template_id`、`belong_member_id`。
  - 只读 Computed 字段：`rule_id`、`sequence`、`source_name`、`dest_name`、`param_template_name`、`edge_name`、`fw_group_name`、`invalid`、`app_id`、`uin`、`nickname`、`member_id`、`node_name`、`create_time`、`update_time`、`dns_parse_cnt`、`hit_count`、`last_hit_time`、`all_total_count`。
- 资源 ID 设为 `GroupId#RuleId`（使用 `tccommon.FILED_SP` 分隔符）。
- Create 路径：调用 `CreateVpcAclRule`，`Rules` 数组固定 1 个元素，取响应 `RuleIds[0]` 作为 RuleId。
- Read 路径：调用 `DescribeVpcAclRules`，用 `GroupId` + `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]` 精准查单条，取 `Rules` 首元素。无单实例 Get 接口。
- Update 路径：调用 `ModifyVpcAclRule`。ModifyVpcAclRule 是全量覆盖（非增量），复用与 Create 相同的字段校验，必须把 state 里所有可写字段完整回传（组 `Rule` 单对象）。不可修改 GroupId。
- Delete 路径：调用 `DeleteVpcAclRule`，`RuleIds` 传 1 个元素。RuleId 已不存在返回 ResourceNotFound 视为删除成功（幂等）。
- 四个接口全部同步（无任务 ID、无中间态），不生成 waitReady/async 逻辑。
- import 支持 `GroupId#RuleId`；Read 遇 ResourceNotFound 时 `d.SetId("")` 移除 state 不报错。
- `source_type`/`dest_type`/`ip_version`/`edge_id`/`fw_group_id` 变更视为重建（forceNew）；其余原地更新。
- 新增 `.md` 资源文档 + 单元测试 `_test.go`（使用 gomonkey mock，不使用 terraform 测试套件）。
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册新资源 `tencentcloud_fwm_vpc_acl_rule`。

## Capabilities

### New Capabilities

- `fwm-vpc-acl-rule-resource`: 新增 `tencentcloud_fwm_vpc_acl_rule` 资源的 schema、CRUD 行为、复合 ID 约定（`GroupId#RuleId`）、字段约束（ForceNew 集合、全量覆盖 Update）、文档与测试规范，作为 FWM VPC 边界 ACL 规则的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go/provider.md 增加注册一行 -->

## Impact

- 新文件：
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`（fwm 全新服务目录，含 `NewFwmService`、`FwmService` 结构、`DescribeFwmVpcAclRuleById`、`DeleteFwmVpcAclRuleById` 方法）
  - `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_vpc_acl_rule_test.go`
  - `website/docs/r/fwm_vpc_acl_rule.html.markdown`（由 `make doc` 生成，收尾阶段处理）
- 既有文件改动：
  - `tencentcloud/connectivity/client.go`：新增 `fwmv20250611Conn` 字段、`UseFwmV20250611Client()` 访问器、`fwm/v20250611` import
  - `tencentcloud/provider.go`：新增 `fwm` 包 import 与资源注册映射 `"tencentcloud_fwm_vpc_acl_rule": fwm.ResourceTencentCloudFwmVpcAclRule()`
  - `tencentcloud/provider.md`：新增资源文档条目（收尾阶段由 `make doc` 生成）
- 依赖：SDK `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（v1.3.196）已通过 go-module-vendor step 下载到 vendor，tools.go blank import 已添加。
- 不修改任何既有资源的 schema，对现有用户零影响。