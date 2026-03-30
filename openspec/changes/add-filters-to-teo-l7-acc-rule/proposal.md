## Why

DescribeL7AccRules API 支持 Filters 参数用于精确筛选查询结果，但当前的 tencentcloud_teo_l7_acc_rule 资源尚未接入此参数。这导致用户无法通过 Filters 参数进行高效的规则筛选，影响使用体验。

## What Changes

- 为 tencentcloud_teo_l7_acc_rule 资源的 read 操作添加 Filters 参数支持
- 在 schema 中定义 Filters 参数结构
- 实现从用户配置到 API 调用的参数映射
- 添加相关测试用例和文档

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-filters-param`: 支持在 tencentcloud_teo_l7_acc_rule 资源中使用 Filters 参数进行 read 操作筛选

### Modified Capabilities

## Impact

- 主要影响 `tencentcloud/services/teo/` 目录下的相关文件
- 需要修改数据源 schema 定义和 API 调用逻辑
- 可能需要更新相关测试用例和文档
