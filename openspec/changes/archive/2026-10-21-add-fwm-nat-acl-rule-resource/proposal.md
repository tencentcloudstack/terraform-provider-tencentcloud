## Why

腾讯云防火墙（FWM）产品的 NAT 边界 ACL 规则当前只能通过控制台或 API 单条管理，Terraform Provider 尚未提供对该规则的声明式管理能力。用户无法用 IaC 编排 NAT 边界 ACL 规则（增删改查），也无法将其与规则组（由独立资源创建）进行关联编排。新增 `tencentcloud_fwm_nat_acl_rule` 资源填补这一空白，使 NAT 边界 ACL 规则全生命周期可被 Terraform 管理。

## What Changes

- 新增资源 `tencentcloud_fwm_nat_acl_rule`（RESOURCE_KIND_GENERAL），基于 `fwm/v20250611` SDK 实现单条 NAT 边界 ACL 规则的 CRUD。
- 接口映射（全部同步，无任务 ID / 无轮询）：
  - 创建：`CreateNatAclRule`（Rules 固定 1 个元素，取返回 RuleIds 首元素绑定 RuleId）。
  - 查询：`DescribeNatAclRules`（Direction 必填；精准查单条用 GroupId + Direction + Filters=[{Name:"RuleId", Values:[ruleId], OperatorType:1}]；取 Rules 首元素）。
  - 更新：`ModifyNatAclRule`（单条全量覆盖更新，需带所有可写字段，不可改 GroupId）。
  - 删除：`DeleteNatAclRule`（RuleIds 传 1 个元素；规则不存在视为删除成功，幂等）。
- 资源 ID 为复合结构 `GroupId + FILED_SP + RuleId`（separator = `#`）。
- 新增 fwm 服务的 connectivity client 绑定 `UseFwmV20250611Client()`（当前 provider 仅存在 vendor SDK 包，未接入 client.go / provider.go）。
- 新增 fwm service 服务层 `service_tencentcloud_fwm.go`，含 `DescribeFwmNatAclRuleById` 查询单条规则。
- Schema 关键字段错位映射（显式处理，避免 import 后 plan 永久 diff）：
  - 入参 `Rule.OrderIndex` ↔ 出参 `Rules[].Sequence`（回读用 Sequence 填 order_index）。
  - 入参 `Rule.BelongMemberId`(string) ↔ 出参 `Rules[].BelongMember`(对象)（回读取 BelongMember.MemberId，仅非空才写入 belong_member_id）。
- ForceNew 策略：`direction`、`source_type`、`target_type` 变更视为重建；`group_id` required + ForceNew；其余规则属性（`source_content`、`target_content`、`protocol`、`rule_action`、`order_index`、`scope`、`port`、`description`、`param_template_id`、`belong_member_id`）原地更新。
- 只读 computed 字段（仅 Describe 回填）：`rule_id`、`sequence`、`scope_desc`、`create_time`、`update_time`、`source_name`、`target_name`、`param_template_name`、`fw_ins_id`、`country_name`、`city_name`、`country_code`、`city_code`、`country_key`、`city_key`、`is_region`、`cloud_code`、`is_cloud`、`invalid`、`app_id`、`uin`、`nickname`、`member_id`、`node_name`、`dns_parse_cnt`、`hit_count`、`last_hit_time`、`all_total_count`。
- import 支持 `GroupId:RuleId` 两段式（实际状态文件存储用 `#` 分隔符）。
- 在 `tencentcloud/provider.go` 注册新资源；生成 `.md` 文档与 `_test.go` 单元测试。

## Capabilities

### New Capabilities
- `fwm-nat-acl-rule-resource`: 管理 FWM 产品下单条 NAT 边界 ACL 规则的全生命周期（创建 / 查询 / 更新 / 删除 / 导入），含关键字段错位映射、同步接口处理、复合 ID 与幂等删除。

### Modified Capabilities
<!-- 本次为纯新增资源，不修改任何已有 spec 级别行为。 -->

## Impact

- **新增代码**:
  - `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.go`：CRUD + schema（单文件，风格参照 `tencentcloud_igtm_strategy`）。
  - `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule.md`：资源文档 + import 示例。
  - `tencentcloud/services/fwm/resource_tc_fwm_nat_acl_rule_test.go`：单元测试（使用 gomonkey mock 云 API，不使用 terraform 测试套件）。
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`：fwm service 层，含 `DescribeFwmNatAclRuleById`。
- **修改代码**:
  - `tencentcloud/connectivity/client.go`：新增 fwm SDK import、`fwmv20250611Conn` 结构体字段、`UseFwmV20250611Client()` 方法。
  - `tencentcloud/provider.go`：新增 `services/fwm` import，在 `ResourcesMap` 注册 `tencentcloud_fwm_nat_acl_rule`。
- **APIs consumed**: `CreateNatAclRule`、`DescribeNatAclRules`、`ModifyNatAclRule`、`DeleteNatAclRule`（均已存在于 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611/`）。
- **无破坏性变更**：纯新增，不修改任何已有资源 schema 或 state。
- **无需 SDK 升级**：所需 API 均已在 vendor SDK 中。