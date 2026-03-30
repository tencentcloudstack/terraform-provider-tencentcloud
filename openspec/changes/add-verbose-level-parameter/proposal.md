## Why

用户在使用 `tencentcloud_nats` 数据源时，需要通过 `VerboseLevel` 参数控制 API 返回数据的详细程度，以便在不同的使用场景下获取合适粒度的信息（如基础信息 vs 详细配置信息），提高查询效率和灵活性。

## What Changes

- 在 `tencentcloud_nats` 数据源的 schema 中新增 `VerboseLevel` 可选参数
- 修改数据源查询逻辑，将 `VerboseLevel` 参数传递给 `DescribeNatGateways` API 调用
- 更新数据源文档，说明 `VerboseLevel` 参数的用法和可选值

## Capabilities

### New Capabilities
- `verbose-level-filter`: 支持 tencentcloud_nats 数据源通过 VerboseLevel 参数控制 API 返回数据的详细程度

### Modified Capabilities
- (无，仅新增参数，不修改现有能力)

## Impact

- 受影响的代码：`tencentcloud/services/vpc/data_source_tc_nats.go`
- 受影响的文档：`website/docs/data-sources/tencentcloud_nats.html.markdown`
- 受影响的测试：`tencentcloud/services/vpc/data_source_tc_nats_test.go`
- 依赖的 API：TencentCloud API DescribeNatGateways（该 API 已支持 VerboseLevel 参数）
