## ADDED Requirements

### Requirement: Set verbose level parameter
用户 SHALL 能够在 `tencentcloud_nats` 数据源配置中设置 `VerboseLevel` 参数，以控制 API 返回数据的详细程度。

#### Scenario: User sets verbose level
- **WHEN** 用户在 `tencentcloud_nats` 数据源配置中设置 `VerboseLevel` 为整数值
- **THEN** 该参数 SHALL 被包含在数据源配置中
- **AND** 该参数 SHALL 被传递给 DescribeNatGateways API 调用

#### Scenario: User does not set verbose level
- **WHEN** 用户在 `tencentcloud_nats` 数据源配置中不设置 `VerboseLevel` 参数
- **THEN** 数据源 SHALL 正常工作
- **AND** DescribeNatGateways API SHALL 使用默认行为

### Requirement: Pass verbose level to API
Provider SHALL 将用户配置的 `VerboseLevel` 参数值正确传递给 TencentCloud DescribeNatGateways API。

#### Scenario: API call includes verbose level
- **WHEN** 用户设置了 `VerboseLevel` 参数
- **THEN** API 调用 SHALL 包含该参数
- **AND** API 返回的数据 SHALL 对应用户指定的详细程度

#### Scenario: API call without verbose level
- **WHEN** 用户未设置 `VerboseLevel` 参数
- **THEN** API 调用 SHALL 不包含该参数
- **AND** API SHALL 返回默认详细程度的数据

### Requirement: Maintain backward compatibility
新增的 `VerboseLevel` 参数 SHALL 不影响现有用户配置和数据源行为。

#### Scenario: Existing configuration works
- **WHEN** 用户的现有配置不包含 `VerboseLevel` 参数
- **THEN** 配置 SHALL 继续正常工作
- **AND** 数据源 SHALL 返回相同的结果
- **AND** 不需要用户修改任何配置
