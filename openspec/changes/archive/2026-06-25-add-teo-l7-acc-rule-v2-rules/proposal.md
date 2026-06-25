## Why

TEO 的 `tencentcloud_teo_l7_acc_rule_v2` 资源当前仅将云API `RuleEngineItem` 的子字段（`status`、`rule_name`、`description`、`branches` 等）平铺在 Terraform schema 中，缺乏一个直接从 API 返回的 `Rules` 结构化输出的顶层字段。新增 `Rules` 参数可以让用户通过单一字段获取完整的规则数据结构，便于与上游系统对接时直接传递原始规则内容。

## What Changes

- 为 `tencentcloud_teo_l7_acc_rule_v2` 资源新增一个名为 `rules` 的参数：
  - 该参数映射 `CreateL7AccRules` 接口的 `request.Rules`（入参）
  - 该参数映射 `DescribeL7AccRules` 接口的 `response.Rules`（出参）
  - 类型为 `schema.TypeList`，元素为 `RuleEngineItem` 结构体

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-rules-output`: 为 TEO L7 加速规则 v2 资源新增 `rules` 参数，支持在 Terraform state 中以结构化形式存储和查看完整的规则数据

### Modified Capabilities
<!-- No existing capabilities are being modified at spec level -->

## Impact

- **代码变更**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go` — 新增 schema 定义、Create/Read 逻辑中补充 `rules` 字段处理
- **文档**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md` — 新增 `rules` 参数的使用示例
- **测试**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2_test.go` — 补充单元测试
- **API 依赖**: 依赖已有的 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包，无需升级
- **向后兼容**: 仅新增 Optional/Computed 参数，不影响已有配置和 state
