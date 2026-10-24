## Why

腾讯云云防火墙（CFW / 产品代号 `fwm`）提供互联网边界（Edge）ACL 规则能力，SDK `fwm/v20250611` 已交付 `CreateEdgeAclRule`、`DescribeEdgeAclRules`、`ModifyEdgeAclRule`、`DeleteEdgeAclRule` 四个同步接口。当前 Terraform Provider 未暴露任何 fwm 资源，用户无法以基础设施即代码的方式管理 Edge ACL 单条规则，从而无法实现规则的版本化、审计与自动化治理。本变更新增通用 CRUD 资源 `tencentcloud_fwm_edge_acl_rule` 以补齐该能力。

## What Changes

- 新增通用（RESOURCE_KIND_GENERAL）资源 `tencentcloud_fwm_edge_acl_rule`，位于新服务包 `tencentcloud/services/fwm/`，文件 `resource_tc_fwm_edge_acl_rule.go`，覆盖单条 Edge ACL 规则的完整生命周期（Create/Read/Update/Delete + Import）。
- 复合资源 ID：`GroupId#RuleId`（使用 `tccommon.FILED_SP` 分隔），创建时从 `CreateEdgeAclRule` 返回的 `RuleIds[0]` 取 `RuleId`；Read/Update/Delete 中拆分使用。
- Schema 字段一对一映射 `CreateEdgeAclRule`/`ModifyEdgeAclRule` 的 `Rule`/`Rules` 字段，并显式处理关键字段错位：入参 `OrderIndex` ↔ 出参 `Sequence`；入参 `BelongMemberId` ↔ 出参 `BelongMember.MemberId`。
- `group_id` 标 Required；`direction` 标 Required + ForceNew；`source_type`/`target_type` 标 ForceNew；`order_index` Required 且可变。其余规则属性原地更新。
- 四个接口均为同步，Create/Modify/Delete 成功即生效，不生成 waitReady/异步轮询逻辑。
- Describe 走列表查询 + `Filters=[{Name:"RuleId",Values:[ruleId],OperatorType:1}]` 精准定位，`Direction` 必填，`Limit` 取 API 注释最大值 1000。
- 只读字段（rule_id、sequence、create_time、update_time 等）仅 Describe 回填，用户不可写。
- 在 `tencentcloud/connectivity/client.go` 新增 `UseFwmV20250611Client()` 客户端访问器（lazy-init，参照 `UseIgtmV20231024Client`）。
- 在 `tencentcloud/provider.go` 注册资源并在 `tencentcloud/provider.md` 添加资源条目。
- 生成 `resource_tc_fwm_edge_acl_rule.md` 文档与 `resource_tc_fwm_edge_acl_rule_test.go` 单元测试（使用 gomonkey mock 云 API，不使用 terraform 测试套件）。

## Capabilities

### New Capabilities
- `fwm-edge-acl-rule-resource`: 管理一条腾讯云云防火墙互联网边界（Edge）ACL 规则的完整生命周期，包括基于 `GroupId#RuleId` 的导入、字段错位的读写映射、以及 `ResourceNotFound` 时幂等删除。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- 新增文件:
  - `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.go`
  - `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule.md`
  - `tencentcloud/services/fwm/resource_tc_fwm_edge_acl_rule_test.go`
  - `tencentcloud/services/fwm/service_tencentcloud_fwm.go`
- 修改文件:
  - `tencentcloud/connectivity/client.go` — 新增 import 别名 `fwmv20250611`、私有缓存字段 `fwmv20250611Conn`、访问器方法 `UseFwmV20250611Client()`
  - `tencentcloud/provider.go` — 新增 `services/fwm` 包 import 与资源注册 `"tencentcloud_fwm_edge_acl_rule": fwm.ResourceTencentCloudFwmEdgeAclRule()`
  - `tencentcloud/provider.md` — 新增资源条目
- SDK 依赖: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611`（namespace `v20250611`），已存在于 `vendor/` 目录。
- 向后兼容: 纯新增资源，不修改任何已有 schema 或 state，完全向后兼容。