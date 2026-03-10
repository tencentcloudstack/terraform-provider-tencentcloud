## ADDED Requirements

### Requirement: Custom Timeout Configuration
`tencentcloud_clb_instance` 资源 MUST 支持 Terraform `timeouts` 块，允许用户自定义 Create 和 Update 操作的超时时间。

#### Scenario: 使用默认超时创建 CLB 实例
- **WHEN** 用户未配置 `timeouts` 块创建 CLB 实例
- **THEN** Create 操作使用默认超时时间（10 分钟）

#### Scenario: 使用自定义超时创建 CLB 实例
- **WHEN** 用户配置 `timeouts { create = "20m" }` 创建 CLB 实例
- **THEN** Create 操作使用用户指定的 20 分钟超时

#### Scenario: 使用自定义超时更新 CLB 实例
- **WHEN** 用户配置 `timeouts { update = "15m" }` 更新 CLB 实例属性
- **THEN** Update 操作中所有异步任务等待使用用户指定的 15 分钟超时

## ADDED Requirements

### Requirement: Create Operation Timeout Handling
`resourceTencentCloudClbInstanceCreate` 中所有异步任务等待（包括创建 CLB、设置安全组、设置日志、修改 target_region_info、设置 delete_protect、关联 endpoint）MUST 使用 `d.Timeout(schema.TimeoutCreate)` 替代硬编码超时。

#### Scenario: Create 中多个异步任务均使用可配置超时
- **WHEN** 创建 CLB 实例并设置安全组、日志、target_region_info 等附加配置
- **THEN** 每个异步任务等待均使用 `d.Timeout(schema.TimeoutCreate)` 作为超时上限

### Requirement: Update Operation Timeout Handling
`resourceTencentCloudClbInstanceUpdate` 中所有异步任务等待（包括修改 SLA、修改属性、设置安全组、设置日志、修改 project、EIP 操作）MUST 使用 `d.Timeout(schema.TimeoutUpdate)` 替代硬编码超时。

#### Scenario: Update 中多个异步任务均使用可配置超时
- **WHEN** 更新 CLB 实例的多个属性触发多个异步任务
- **THEN** 每个异步任务等待均使用 `d.Timeout(schema.TimeoutUpdate)` 作为超时上限
