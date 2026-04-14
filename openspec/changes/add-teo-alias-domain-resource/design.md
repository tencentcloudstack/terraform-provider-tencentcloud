## Context

Terraform Provider for TencentCloud 目前缺少对 Teo 产品的别称域名（Alias Domain）资源支持。Teo 是腾讯云的边缘安全加速产品，别称域名功能允许用户为站点配置别名域名。当前用户需要通过云 API 或控制台手动管理别称域名，无法通过 Terraform 进行声明式管理。

## Goals / Non-Goals

**Goals:**
- 实现 Teo 别称域名的完整 CRUD 操作
- 支持通过 Terraform 声明式管理别称域名的生命周期
- 确保异步操作的可靠性，通过轮询机制确认操作生效
- 提供完整的测试覆盖和文档

**Non-Goals:**
- 不支持批量创建或删除别称域名（云 API 只支持单个操作）
- 不支持别称域名的其他高级配置（如自定义规则等，仅支持基本配置）
- 不涉及 Teo 其他资源的修改或新增

## Decisions

### 1. 异步操作的轮询机制
**决策**: 对所有标记为异步的云 API（CreateAliasDomain、ModifyAliasDomain、ModifyAliasDomainStatus、DeleteAliasDomain），在调用后立即调用 DescribeAliasDomains 接口轮询，直到资源状态与预期一致或达到超时。

**理由**:
- Teo 云 API 的创建、修改、删除、状态修改都是异步操作，立即返回可能表示操作已受理但未实际生效
- Terraform 期望资源状态在 apply 完成后是最终一致的状态，轮询机制可以保证这一点
- 使用 `helper.Retry()` 函数实现轮询，与 provider 现有模式保持一致

**备选方案**:
- 使用 state refreshers: 功能类似，但 helper.Retry() 更简单直接
- 直接返回不轮询: 会导致 Terraform state 与实际状态不一致，不符合最佳实践

### 2. 资源 ID 格式
**决策**: 使用复合 ID 格式 `zone_id#alias_name`，其中 `zone_id` 是站点 ID，`alias_name` 是别称域名名称。

**理由**:
- 别称域名由 `zone_id` 和 `alias_name` 唯一标识
- 使用 `#` 分隔符与 provider 现有模式一致（如 `instanceId#userId`）
- 在 Read 操作中可以从 ID 解析出这两个关键字段

### 3. 状态暂停/启用机制
**决策**: 通过 `paused` 参数控制别称域名的暂停/启用状态，true 表示暂停，false 表示启用。在 Update 函数中检测 `paused` 参数的变化，如果变化则调用 ModifyAliasDomainStatus API。

**理由**:
- 云 API 的 ModifyAliasDomainStatus 接口专门用于修改状态
- 将状态修改集成到 Update 函数中，用户可以通过修改配置文件来切换状态
- 需要与其他参数的修改（通过 ModifyAliasDomain）分开处理

### 4. Timeout 配置
**决策**: 在 Schema 中声明 Timeouts 块，提供 Create、Update、Delete 的自定义超时配置。

**理由**:
- 异步操作可能需要较长时间才能生效，默认超时可能不够
- 允许用户根据实际情况调整超时时间，提高配置的灵活性
- 符合 Terraform provider 的最佳实践

## Risks / Trade-offs

### 1. 轮询超时导致的状态不一致
**风险**: 如果云 API 操作失败或网络问题导致轮询超时，Terraform state 可能与实际状态不一致。

**缓解措施**:
- 使用合理的默认超时时间（建议 10 分钟）
- 在轮询失败时返回明确的错误信息
- 提供文档说明，指导用户在超时后使用 `terraform refresh` 同步状态

### 2. 云 API 行为变化
**风险**: 云 API 的行为或参数可能在未来版本中变化，导致 provider 代码需要更新。

**缓解措施**:
- 使用 vendor 模式管理依赖，锁定具体版本
- 在测试中覆盖主要场景，尽早发现 API 变化
- 添加适当的错误处理和日志记录

### 3. 并发操作冲突
**风险**: 同一个别称域名被多个 Terraform 配置并发操作，可能导致状态不一致。

**缓解措施**:
- 云 API 本身应该有并发控制，依赖其实现
- 提示用户避免对同一资源进行并发操作
- 在文档中说明最佳实践

## Migration Plan

由于这是一个全新的资源，不涉及现有资源的迁移或数据迁移。

**部署步骤**:
1. 代码审查和测试
2. 发布新版本 provider
3. 用户更新 provider 版本后即可使用新资源

**回滚策略**:
- 如果新资源有严重问题，可以发布新版本移除该资源
- 不会影响现有资源和配置的向后兼容性

## Open Questions

目前没有未决的问题。所有技术决策都已经明确。
