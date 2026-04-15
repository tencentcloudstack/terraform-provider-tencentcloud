## Context

当前 Terraform Provider for TencentCloud 已经实现了多个云产品的资源管理，但 Teo 产品（边缘安全加速）的 DNS 记录管理能力缺失。Teo 是腾讯云的边缘安全加速产品，提供 DNS 记录管理、边缘加速等功能。

项目的代码组织结构遵循以下约定：
- 资源实现文件位于 `tencentcloud/services/<service>/resource_tc_<service>_<name>.go`
- 使用 `tencentcloud-sdk-go` 调用云 API
- 遵循 Terraform Plugin SDK v2 的资源开发模式
- 所有异步操作必须支持 Timeouts 配置
- 使用 helper.Retry() 实现最终一致性重试
- 错误处理使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`

Teo 产品的 API 位于 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 包中，本次变更需要使用该包的四个接口来管理 DNS 记录。

## Goals / Non-Goals

**Goals:**
- 实现完整的 DNS 记录 CRUD 操作，满足用户的资源管理需求
- 遵循项目现有的代码规范和模式，保持代码一致性
- 支持异步操作的超时配置，确保资源操作的可靠性
- 提供完善的单元测试，确保代码质量
- 提供清晰的文档和使用示例，降低用户使用门槛

**Non-Goals:**
- 不涉及 Teo 产品的其他资源（如站点配置、边缘规则等）
- 不涉及现有资源的修改或迁移
- 不实现复杂的批量操作（虽然底层 API 支持批量，但 Terraform 资源层面每个资源实例对应一条 DNS 记录）

## Decisions

### 1. 资源 ID 设计

**决策:** 采用复合 ID 结构，格式为 `zone_id#record_id`

**理由:**
- DNS 记录需要 Zone ID 和记录 ID 来唯一标识
- 使用 `#` 作为分隔符是项目中的标准做法
- 复合 ID 便于在 CRUD 操作中解析和使用
- 遵循项目现有的复合 ID 约定

**替代方案考虑:**
- 单一 ID: 无法唯一标识记录，因为同一个 Zone 内可能有相同的记录 ID
- 使用逗号分隔: 项目中统一使用 `#` 分隔符，保持一致性

### 2. 异步操作处理

**决策:** 为异步 API 调用实现 Read 接口轮询机制

**理由:**
- Terraform 资源的状态需要与云端保持一致
- 云 API 的异步操作需要时间生效，需要轮询确认
- 在 schema 中声明 Timeouts 块，允许用户配置超时时间
- 使用 ctx 传递超时上下文，避免无限等待

**实现方式:**
- 在 Create/Update/Delete 操作后，调用 Read 接口轮询
- 使用 helper.Retry() 或自定义重试逻辑
- 达到超时时间后返回错误

### 3. Schema 设计原则

**决策:** 严格遵循云 API 的参数定义，只映射云 API 支持的参数

**理由:**
- 确保资源行为与云 API 完全一致
- 避免参数不匹配导致的运行时错误
- 简化维护成本，当云 API 更新时可以同步更新

**实现方式:**
- 根据 vendor 目录下的云 API 定义来设计 schema
- 对于 Create 接口，只包含云 API CreateDnsRecord 支持的参数
- 对于 Update 接口，只包含云 API ModifyDnsRecords 支持的可修改参数
- 使用正确的类型映射（String, Int, Bool, List, Map 等）
- 设置合理的默认值和约束

### 4. 服务层封装

**决策:** 在 service_tencentcloud_teo.go 中封装 API 调用逻辑

**理由:**
- 分离资源层和服务层，提高代码可测试性
- 便于复用 API 调用逻辑（如参数转换、错误处理）
- 遵循项目现有的分层架构

**实现方式:**
- 创建 `CreateDnsRecord`, `DescribeDnsRecords`, `ModifyDnsRecords`, `DeleteDnsRecords` 服务层方法
- 处理 API 请求和响应的参数转换
- 统一处理 API 错误

### 5. 错误处理和重试

**决策:** 使用项目的标准错误处理模式

**理由:**
- 保持代码一致性
- 利用项目现有的错误处理工具
- 确保最终一致性和可靠性

**实现方式:**
- 使用 `defer tccommon.LogElapsed()` 记录操作耗时
- 使用 `defer tccommon.InconsistentCheck()` 检查状态一致性
- 对于可重试的错误，使用 `helper.Retry()` 实现重试逻辑

## Risks / Trade-offs

### 风险 1: 云 API 接口变更或不稳定

**风险:** Teo 产品的 API 可能会变更，或者接口在某些情况下不稳定。

**缓解措施:**
- 使用 vendor 模式管理依赖，锁定特定版本的 SDK
- 在服务层增加完善的错误处理和重试机制
- 在单元测试中模拟各种 API 响应场景
- 提供详细的错误信息，便于用户排查问题

### 风险 2: 异步操作轮询超时

**风险:** 如果云端的异步操作耗时较长，可能导致 Terraform 操作超时失败。

**缓解措施:**
- 在 schema 中声明 Timeouts 块，允许用户配置超时时间
- 使用合理的默认超时时间（如 10 分钟）
- 在超时时返回清晰的错误信息，提示用户调整超时配置

### 风险 3: 复合 ID 解析错误

**风险:** 如果复合 ID 格式不正确，可能导致解析失败。

**缓解措施:**
- 在代码中增加 ID 格式校验
- 提供清晰的错误信息，指出正确的 ID 格式
- 在单元测试中覆盖各种 ID 格式场景

### 权衡 1: 批量操作 vs 单条操作

**权衡:** 底层的 ModifyDnsRecords 和 DeleteDnsRecords API 支持批量操作，但 Terraform 资源层面采用单条操作。

**决策:** 采用单条操作模式，每个资源实例对应一条 DNS 记录。

**理由:**
- 符合 Terraform 资源的语义模型
- 用户可以通过创建多个资源实例来管理多条记录
- 简化实现和维护成本
- 批量操作可以在资源层面通过多个资源实例实现

### 权衡 2: 服务层封装程度

**权衡:** 服务层封装程度过高可能影响灵活性，过低可能无法提供足够的抽象。

**决策:** 保持适度的封装，只封装必要的 API 调用和参数转换逻辑。

**理由:**
- 避免过度设计，保持代码简洁
- 资源层仍然可以直接访问服务层方法
- 便于后续扩展和维护
