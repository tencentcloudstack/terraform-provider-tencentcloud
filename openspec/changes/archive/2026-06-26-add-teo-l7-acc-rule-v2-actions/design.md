## Context

当前 `tencentcloud_teo_l7_acc_rule_v2` 资源（文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`）通过 `TencentTeoL7RuleBranchBasicInfo(1)` 在 `branches` 子 schema 中已包含 `actions` 字段。该字段对应 SDK `RuleBranch.Actions`（类型：`[]*RuleEngineAction`），用于定义加速规则的操作列表。

SDK 中 `RuleEngineAction` 结构体已包含以下三个字段但 Terraform 资源尚未支持：
- `AdvancedOriginRoutingParameters *AdvancedOriginRoutingParameters`：高级回源优化配置
- `ShieldParameters *ShieldParameters`：源站卸载（Shield）配置
- `SiteFailoverParameters *SiteFailoverParameters`：源站故障转移配置

SDK 结构体详情：
- `AdvancedOriginRoutingParameters`：包含 `Direction *string` 字段，取值为 `MainlandChinaAndGlobalAdaptive`
- `ShieldParameters`：包含 `ShieldSpaceId *string` 字段
- `SiteFailoverParameters`：包含 `SiteFailoverStatusCodes []*int64`（取值 4xx/5xx）和 `SiteFailoverParams []*SiteFailover` 字段
- `SiteFailover`：包含 `Mode`、`Origin`、`OriginProtocol`、`HTTPOriginPort`、`HTTPSOriginPort`、`UpstreamHostHeader`、`UpstreamURLRewrite`、`UpstreamRequestParameters`、`UpstreamHTTP2Parameters`、`PrivateAccess`、`PrivateParameters`、`RedirectURL`、`ResponsePageId`、`StatusCode` 等字段

## Goals / Non-Goals

**Goals:**
- 在 `branches.actions` 的 Schema 中新增 `advanced_origin_routing_parameters` 参数，类型为 TypeList，MaxItems: 1，元素包含 `direction` 字段
- 在 `branches.actions` 的 Schema 中新增 `shield_parameters` 参数，类型为 TypeList，MaxItems: 1，元素包含 `shield_space_id` 字段
- 在 `branches.actions` 的 Schema 中新增 `site_failover_parameters` 参数，类型为 TypeList，MaxItems: 1，元素包含 `site_failover_status_codes` 和 `site_failover_params`
- Create 时：将这三个参数映射到 SDK `RuleEngineAction` 对应字段
- Read 时：从 API 响应中提取值并设置到 terraform state
- Update 时：当参数有变更时，映射到 SDK 对应字段进行更新
- 在 `actions.name` 描述中添加 `AdvancedOriginRouting`、`Shield`、`SiteFailover` 枚举值
- 保持向后兼容：三个参数均为 Optional，现有配置不受影响

**Non-Goals:**
- 不在资源顶层新增 `actions` 参数（actions 已在 `branches` 中支持，之前错误添加的顶层 `actions` 参数需回退）
- 不修改 `DeleteL7AccRules` 接口调用（该接口不涉及 Actions）
- 修改 `resource_tc_teo_l7_acc_rule_v2.go` 仅限于回退之前错误添加的顶层 `actions` 参数

## Decisions

### 1. 参数定义在 `branches.actions` 内部
**选择**：将三个新参数作为 `TencentTeoL7RuleBranchBasicInfo` 中 `actions` 子 schema 的字段添加。

**理由**：与 SDK 结构一致，`AdvancedOriginRoutingParameters`、`ShieldParameters`、`SiteFailoverParameters` 都属于 `RuleEngineAction` 的字段。

### 2. SiteFailover 嵌套结构完整映射
**选择**：完整映射 `SiteFailover` 结构的所有字段到 terraform schema。

**理由**：`SiteFailover` 结构较复杂但各字段都有实际用途，完整映射可确保用户能配置所有选项。

### 3. Read 时 nil 检查
**选择**：在 Read 的 flatten 逻辑中，对 `AdvancedOriginRoutingParameters`、`ShieldParameters`、`SiteFailoverParameters` 进行 nil 检查后再提取字段值。

**理由**：遵循项目规范，避免 nil 指针异常。

### 4. 回退顶层 `actions` 参数
**选择**：移除之前在 `resource_tc_teo_l7_acc_rule_v2.go` 中错误添加的顶层 `actions` 参数及其 CRUD 逻辑。

**理由**：`actions` 字段已在 `branches` 中支持，无需在资源顶层重复添加。顶层 `actions` 参数是之前迭代中的误操作，需要回退以保持资源接口的简洁性。

## Risks / Trade-offs

- **[风险] 参数数量增加**：`site_failover_parameters` 内部嵌套较深，可能增加用户配置复杂度。
  - **缓解**：参数均为 Optional，不影响已有配置。

## Migration Plan

1. 在 `TencentTeoL7RuleBranchBasicInfo` 的 `actions` schema 中新增三个参数定义
2. 在 `resourceTencentCloudTeoL7AccRuleGetBranchs` 中添加三个参数的 flatten 逻辑
3. 在 `resourceTencentCloudTeoL7AccRuleSetBranchs` 中添加三个参数的 set 逻辑
4. 在 `actions.name` 描述中添加三个枚举值
5. 更新 `.md` 文档添加 Example Usage 展示新参数用法
6. 补充单元测试
