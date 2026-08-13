## Why

为 `tencentcloud_teo_l7_acc_rule_v2` 资源的 `branches.actions` 新增三个操作参数：`advanced_origin_routing_parameters`、`shield_parameters` 和 `site_failover_parameters`。当前资源的 `branches.actions` 已包含多个操作参数（如 Cache、CacheKey 等），但缺少 SDK `RuleEngineAction` 中已有的 `AdvancedOriginRoutingParameters`、`ShieldParameters` 和 `SiteFailoverParameters` 字段的支持。

## What Changes

- 在 `tencentcloud_teo_l7_acc_rule_v2` 资源 Schema 的 `branches.actions` 中新增 Optional 参数 `advanced_origin_routing_parameters`（类型：TypeList，MaxItems: 1），对应 SDK `RuleEngineAction.AdvancedOriginRoutingParameters`，字段包含 `direction`
- 在 `branches.actions` 中新增 Optional 参数 `shield_parameters`（类型：TypeList，MaxItems: 1），对应 SDK `RuleEngineAction.ShieldParameters`，字段包含 `shield_space_id`
- 在 `branches.actions` 中新增 Optional 参数 `site_failover_parameters`（类型：TypeList，MaxItems: 1），对应 SDK `RuleEngineAction.SiteFailoverParameters`，字段包含 `site_failover_status_codes` 和 `site_failover_params`（嵌套 `SiteFailover` 结构）
- 在 Create / Read / Update 方法中添加这三个参数的读写逻辑
- 在 `actions.name` 描述中添加 `AdvancedOriginRouting`、`Shield`、`SiteFailover` 枚举值
- 新增参数均为 Optional，不影响已有配置的向后兼容性

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-actions`: 在 `tencentcloud_teo_l7_acc_rule_v2` 资源的 `branches.actions` 中新增 `advanced_origin_routing_parameters`、`shield_parameters`、`site_failover_parameters` 三个操作参数，支持高级回源优化、源站卸载（Shield）和源站故障转移功能

### Modified Capabilities
（无）

## Impact

- 受影响文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_extension.go`（新增 Schema 定义及 CRUD 逻辑）
- 受影响文件：`tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`（更新文档示例）
- 不受影响：vendor SDK 无需更新（`AdvancedOriginRoutingParameters`、`ShieldParameters`、`SiteFailoverParameters` 字段已存在于 SDK 中）
- 向后兼容：三个参数均为 Optional 参数，不填写时行为不变
