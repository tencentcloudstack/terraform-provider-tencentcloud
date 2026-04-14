## Context

Terraform Provider for TencentCloud 目前已支持 TEO（EdgeOne）服务的多个资源，但缺少对别称域名（Alias Domain）的管理功能。别称域名允许用户为现有的站点域名设置别名，用于流量转发和管理。

TEO 云 API 提供了完整的别称域名管理接口，包括：
- CreateAliasDomain - 创建别称域名
- DescribeAliasDomains - 查询别称域名列表
- ModifyAliasDomain - 修改别称域名目标
- ModifyAliasDomainStatus - 修改别称域名状态（启用/暂停）
- DeleteAliasDomain - 删除别称域名

其中 CreateAliasDomain、ModifyAliasDomain、ModifyAliasDomainStatus 和 DeleteAliasDomain 都是异步接口，需要轮询 DescribeAliasDomains 来确认操作完成。

## Goals / Non-Goals

**Goals:**
- 实现完整的别称域名 CRUD 功能，覆盖创建、读取、更新、删除生命周期
- 正确处理异步操作的轮询等待机制
- 支持通过 paused 参数管理别称域名的启用/暂停状态
- 确保资源状态的幂等性和一致性
- 提供完整的单元测试，使用 mock 云 API 方式验证代码逻辑

**Non-Goals:**
- 不实现别称域名的批量操作（如批量创建、批量删除）
- 不实现别称域名的数据源（DataSource）
- 不支持别称域名的 DNS 记录管理（这是 TEO 的其他功能）

## Decisions

### 1. 资源 ID 格式
**决策**: 使用 `zone_id#alias_name` 作为复合 ID，以 `#` 分隔
**理由**:
- zone_id 是站点级别的标识符，alias_name 是别称域名的唯一标识
- 复合 ID 能够唯一确定一个别称域名资源
- 遵循项目中其他资源的复合 ID 模式（如 `instanceId#userId`）

### 2. 异步操作轮询策略
**决策**: 对于 CreateAliasDomain、ModifyAliasDomain、ModifyAliasDomainStatus、DeleteAliasDomain 等异步接口，使用 `helper.Retry()` 函数轮询 DescribeAliasDomains 直到状态符合预期
**理由**:
- 项目中已存在成熟的最终一致性重试模式
- 通过查询接口确认操作完成，确保 TF 状态与云资源状态一致
- 避免固定的延迟等待，提高效率

### 3. paused 状态的实现方式
**决策**: 将 paused 作为可选的 Computed 字段，通过 ModifyAliasDomainStatus 接口单独管理状态变更
**理由**:
- paused 状态可以通过 ModifyAliasDomainStatus 独立修改，不影响其他字段
- 设置为 Computed 允许用户查询当前状态，同时支持显式设置
- 将状态管理与其他字段修改解耦，减少不必要的 API 调用

### 4. Update 函数的实现策略
**决策**: Update 函数分别处理不同场景：
- 如果 target_name 发生变化，调用 ModifyAliasDomain
- 如果 paused 发生变化，调用 ModifyAliasDomainStatus
- 如果两者都变化，依次调用两个接口
**理由**:
- ModifyAliasDomain 和 ModifyAliasDomainStatus 是两个独立的 API 接口
- 根据实际变化的字段调用对应的接口，避免不必要的操作
- 确保每个异步操作都能正确轮询确认

### 5. Delete 函数的幂等性处理
**决策**: Delete 函数在调用 DescribeAliasDomains 时，如果发现资源已不存在（ResourceNotFound），直接返回成功
**理由**:
- 确保删除操作的幂等性
- 避免因为资源已被外部删除而导致的删除失败

### 6. 单元测试策略
**决策**: 使用 mock 云 API 的方式编写单元测试，不调用真实的云 API
**理由**:
- 单元测试应该快速、可重复，不依赖外部环境
- mock 可以模拟各种场景（成功、失败、超时等）
- 确保代码逻辑的正确性，避免集成测试的复杂性

## Risks / Trade-offs

### 风险 1: 异步操作超时
**描述**: 轮询 DescribeAliasDomains 可能在默认超时时间内未返回预期状态
**缓解措施**:
- 在 schema 中声明 Timeouts 块，允许用户自定义超时时间
- 提供合理的默认超时时间（如 10 分钟）
- 在日志中记录轮询状态，便于问题排查

### 风险 2: 并发修改导致状态不一致
**描述**: 多个 Terraform 操作或外部系统同时修改同一个别称域名可能导致状态冲突
**缓解措施**:
- Update 函数基于读取到的最新状态进行修改
- 使用 Terraform 的状态锁机制避免并发更新
- 在文档中说明资源应避免被外部系统同时修改

### 风险 3: paused 状态与资源存在性不一致
**描述**: 用户可能混淆 paused=false（启用状态）和资源已删除的概念
**缓解措施**:
- 在文档中明确说明 paused 参数的含义和范围
- paused 只影响别称域名的启用状态，不影响资源的存在性
- 提供清晰的错误信息，区分状态变更失败和资源不存在的情况

### 权衡 1: 轮询频率 vs API 调用次数
**描述**: 轮询频率过高会增加 API 调用次数，过低会增加等待时间
**权衡选择**:
- 采用项目中常用的轮询间隔（如 5-10 秒）
- 平衡效率和成本，符合云 API 的最佳实践

### 权衡 2: 状态管理 vs 操作复杂度
**描述**: 将 paused 作为独立状态管理增加了 Update 函数的复杂度
**权衡选择**:
- 保持代码的清晰性和可维护性，通过函数分解处理不同场景
- 提供充分的注释和单元测试，确保逻辑正确
