## Why

DescribeL7AccRules API 返回的数据当前缺少 TotalCount 字段，导致无法获取查询结果的总条数。当用户需要分页查询或获取规则总数时，这个字段非常重要，因为它可以帮助用户了解数据总量，便于进行分页处理和总量统计。

## What Changes

- 在 DescribeL7AccRules API 响应中添加 TotalCount 字段
- 更新相关的 Go 结构体以包含 TotalCount 参数
- 确保 TotalCount 字段在数据源读取时正确返回

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-totalcount`: 添加 TotalCount 参数支持，使数据源能够返回规则总数信息

### Modified Capabilities
- 无现有能力需要修改

## Impact

- 影响范围：tencentcloud/services/teo/ 目录下的相关文件
- 主要影响：
  - service_tencentcloud_teo.go 中的 DescribeTeoL7AccRuleById 函数
  - 可能需要更新的数据源 schema（如果存在）
  - 相关的测试文件
- 无破坏性变更：这是新增字段，不影响现有功能
