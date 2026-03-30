## Why

为 tencentcloud_nats 数据源添加 VerboseLevel 参数支持，允许用户根据需要获取不同详细程度的 NAT 网关信息，优化查询性能并减少不必要的网络传输和数据解析开销。

## What Changes

- 在 `tencentcloud_nats` 数据源的 schema 中新增 `verbose_level` 可选参数
- 支持三种详细级别：
  - `DETAIL`: 输出实例所有信息（包括 NAT 规则和自定义路由）
  - `COMPACT`: 不输出 NAT 规则和自定义路由，输出实例基本信息、特性开关和 EIP 信息
  - `SIMPLE`: 仅输出实例基本信息和特性开关
- 在调用 DescribeNatGateways API 时将 `verbose_level` 参数值传递给 `VerboseLevel` 字段

## Capabilities

### New Capabilities
- `nat-gateway-verboselevel`: 为 NAT 网关数据源添加详细级别查询控制能力

### Modified Capabilities
- 无

## Impact

- **修改的代码**:
  - `tencentcloud/services/vpc/data_source_tc_nats.go`: 在 schema 中新增 verbose_level 参数，在 Read 函数中传递给 API
- **API 变更**:
  - DescribeNatGateways API 调用将新增 VerboseLevel 参数
- **兼容性**:
  - 新增 Optional 参数，不影响现有配置的向后兼容性
  - 若不指定 verbose_level，API 将使用默认行为
- **文档**:
  - 需要更新 `website/docs/d/tencentcloud_nats.html.markdown` 文档
