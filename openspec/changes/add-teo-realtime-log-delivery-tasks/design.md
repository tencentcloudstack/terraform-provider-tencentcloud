## Context

tencentcloud_teo_realtime_log_delivery 资源当前可以创建和管理实时日志推送任务，但在读取时仅返回当前任务的基本信息。DescribeRealtimeLogDeliveryTasks API 已在服务层通过 `DescribeTeoRealtimeLogDeliveryById` 函数实现，该函数可以返回完整的任务信息，包括任务列表及其详细状态。为了提高资源的可观测性，需要将这些详细信息作为只读参数暴露给 Terraform 用户。

当前状态：
- 资源已支持基本的 CRUD 操作
- 服务层已有 DescribeRealtimeLogDeliveryTasks API 的调用实现
- Read 函数已能获取单个任务的基本信息
- 缺少 RealtimeLogDeliveryTasks 参数来展示任务的详细信息

约束：
- 必须保持向后兼容，不能破坏现有 TF 配置和 state
- 只能新增 Optional 字段，不能修改已有资源的 schema
- 需要遵循 Terraform Provider SDK v2 的最佳实践

## Goals / Non-Goals

**Goals:**
- 为 tencentcloud_teo_realtime_log_delivery 资源添加 RealtimeLogDeliveryTasks 只读参数
- 通过 DescribeRealtimeLogDeliveryTasks API 获取并填充该参数
- 确保新增参数与现有 schema 完全兼容
- 提供清晰的参数文档说明

**Non-Goals:**
- 不修改资源的 Create、Update、Delete 操作
- 不改变现有参数的行为
- 不引入新的外部依赖

## Decisions

### 1. Schema 设计
- 在资源 Schema 中添加 `realtime_log_delivery_tasks` 字段，类型为 `TypeList`，属性设为 `Computed` 和 `Optional`
- 该字段包含任务列表的详细信息，每个任务是一个嵌套对象
- 使用现有的 `DescribeTeoRealtimeLogDeliveryById` 服务层函数获取数据

**理由：**
- `TypeList` + `Computed` + `Optional` 是 Terraform 中只读参数的标准模式
- 不会影响现有的配置文件，因为它是可选的
- 复用现有的服务层函数，避免重复代码

**替代方案考虑：**
- 使用数据源（DataSource）：不合适，因为这是为现有资源添加属性，不是创建新的查询接口
- 修改现有参数：会破坏向后兼容性，违反硬约束

### 2. 数据映射
- 将 API 返回的 `RealtimeLogDeliveryTasks` 结构直接映射到 Terraform schema
- 保留 API 返回的所有关键字段，如 TaskId、ZoneId、TaskName、TaskType、Status 等
- 使用嵌套的 `TypeList` + `TypeMap` 结构来表示复杂字段

**理由：**
- 直接映射可以保持数据的完整性
- 减少数据转换的复杂度
- 用户可以直接使用 API 文档了解各字段的含义

### 3. 实现位置
- 修改 `resource_tc_teo_realtime_log_delivery.go` 文件
- 在 `ResourceTencentCloudTeoRealtimeLogDelivery` 函数中添加新字段到 Schema
- 在 `resourceTencentCloudTeoRealtimeLogDeliveryRead` 函数中调用服务层函数并填充新字段

**理由：**
- 集中在单个文件中，便于维护
- 遵循项目现有的代码组织结构
- Read 函数是设置只读参数的正确位置

## Risks / Trade-offs

### 风险 1: API 兼容性
如果 DescribeRealtimeLogDeliveryTasks API 的返回结构发生变化，可能导致数据映射失败。

**缓解措施：**
- 使用 tencentcloud-sdk-go 提供的类型，避免直接处理 JSON
- 在代码中添加错误处理，如果 API 返回数据不符合预期，记录警告但不中断读取流程
- 保持与 SDK 版本同步更新

### 风险 2: 性能影响
每次读取资源时都会调用 API，可能增加响应时间。

**缓解措施：**
- 使用现有的重试机制（helper.Retry）确保最终一致性
- API 调用已存在，只是复用现有数据，不会增加额外的网络请求
- 如果性能成为问题，可以考虑使用缓存机制（当前不需要）

### 权衡 1: 数据完整性 vs 复杂度
保留 API 返回的所有字段可以保证数据完整性，但会增加 schema 的复杂度。

**决策：** 保留所有关键字段，因为：
- 用户可以灵活选择需要的字段
- 不限制 API 的演进空间
- Terraform 的 schema 足够灵活，可以处理复杂结构

### 权衡 2: 只读 vs 可读写
将参数设为只读可以避免误操作，但限制了某些使用场景。

**决策：** 使用只读参数，因为：
- 这些字段由 API 返回，用户不应修改
- Create/Update 操作已有对应的配置参数
- 符合 "读取" 的语义

## Migration Plan

由于这是只新增可选参数，不需要迁移计划：
- 现有配置文件无需修改
- 现有 state 无需迁移
- 新参数对现有用户完全透明
- 用户可以升级到新版本后，按需使用新参数

回滚策略：
- 如果出现问题，可以通过回退 Terraform Provider 版本来移除新参数
- 由于参数是可选的，回退不会影响已部署的资源

## Open Questions

无。这是一个相对简单直接的变更，所有技术决策都已明确。
