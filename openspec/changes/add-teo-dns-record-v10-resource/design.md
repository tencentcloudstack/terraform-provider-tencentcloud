## Context

Terraform Provider for TencentCloud 目前需要为 TEO（EdgeOne）服务新增 DNS 记录管理资源。TEO 服务提供了完整的 DNS 记录管理 API，包括创建、查询、修改和删除操作。

当前 TEO 服务已有多个资源实现，如 `tencentcloud_teo_acceleration_domain` 等，遵循标准的 Terraform Plugin SDK v2 模式。新资源需要与现有资源保持一致的代码风格和实现模式。

TEO DNS 记录 API 具有以下特点：
- CreateDnsRecord：创建单条 DNS 记录，返回 record_id
- DescribeDnsRecords：查询 DNS 记录列表，支持过滤条件
- ModifyDnsRecords：批量修改 DNS 记录，一次最多修改 100 条
- DeleteDnsRecords：批量删除 DNS 记录，一次最多删除 1000 条
- 支持多种 DNS 记录类型：A、AAAA、MX、CNAME、TXT、NS、CAA、SRV
- 部分参数仅用于特定记录类型（如 Location 仅适用于 A、AAAA、CNAME）

## Goals / Non-Goals

**Goals:**

1. 实现完整的 CRUD 操作，支持创建、读取、更新、删除 DNS 记录
2. 支持所有 TEO 支持的 DNS 记录类型（A、AAAA、MX、CNAME、TXT、NS、CAA、SRV）
3. 支持配置解析线路、缓存时间、权重、MX 优先级等高级参数
4. 实现资源状态管理，确保资源创建/更新后生效
5. 提供完整的单元测试和文档

**Non-Goals:**

1. 不支持 DNS 记录的批量导入/导出（通过其他数据源实现）
2. 不支持 DNS 记录的历史版本管理
3. 不支持跨站点的 DNS 记录迁移

## Decisions

### 1. 资源 ID 设计

**决策**: 使用 `zone_id#record_id` 作为资源 ID

**理由**:
- ZoneId 和 RecordId 的组合能唯一标识一个 DNS 记录
- 与现有 TEO 资源保持一致的复合 ID 模式（如 `tencentcloud_teo_acceleration_domain` 使用 `zone_id#domain_name`）
- 便于资源的导入和管理

### 2. Schema 设计

**决策**: 所有云 API 参数都映射为 Terraform Schema 字段，只读字段（如 Status、CreatedOn、ModifiedOn）设置为 Computed

**理由**:
- 完整暴露云 API 能力，用户可以配置所有支持的参数
- 只读字段作为 Computed 字段提供完整的信息反馈
- 符合 Terraform Provider 的最佳实践

### 3. 更新策略

**决策**: 使用 ModifyDnsRecords API 进行更新，需要传递完整的 DnsRecord 结构

**理由**:
- ModifyDnsRecords API 需要完整的记录信息进行更新
- 在 Update 函数中先调用 DescribeDnsRecords 获取当前记录，然后合并变更并调用 ModifyDnsRecords
- 确保更新的原子性和一致性

### 4. 异步操作处理

**决策**: 在 Create 和 Update 操作后，调用 DescribeDnsRecords 进行轮询，直到记录生效

**理由**:
- DNS 记录的创建和修改是异步操作，需要等待生效
- 通过轮询确保状态一致性
- 使用 resource.Retry 实现带超时的重试机制

### 5. 时间配置

**决策**: 在 Schema 中声明 Timeouts 块，设置合理的默认值

**理由**:
- DNS 记录的创建/修改可能需要较长时间生效
- 提供可配置的超时时间，避免长时间等待
- 与现有 TEO 资源保持一致（Create/Update: 20 分钟，Read: 3 分钟，Delete: 20 分钟）

### 6. 错误处理

**决策**: 使用 tccommon.RetryError 进行可重试错误的处理，使用 defer tccommon.LogElapsed() 和 defer tccommon.InconsistentCheck() 进行日志和一致性检查

**理由**:
- 符合现有 TEO 资源的错误处理模式
- 提供完整的错误日志和一致性保证
- 便于问题排查和调试

### 7. 服务层设计

**决策**: 在 service_tencentcloud_teo.go 中添加 DescribeDnsRecordById 方法，用于通过 ID 查询单条记录

**理由**:
- DescribeDnsRecords API 返回列表，需要通过 Filter 参数过滤
- 封装查询逻辑，简化资源层的代码
- 提高代码的可维护性

### 8. 记录类型参数验证

**决策**: 不在 Schema 层面进行参数验证，依赖云 API 的参数校验

**理由**:
- 不同记录类型对参数有不同的要求（如 Location 仅适用于 A、AAAA、CNAME）
- 云 API 会进行参数校验并返回明确的错误信息
- 避免在 Provider 层面维护复杂的校验逻辑

## Risks / Trade-offs

### 风险 1: 异步操作超时

**风险**: DNS 记录的创建/修改可能因网络问题或服务负载导致长时间未生效，超过超时时间

**缓解措施**:
- 设置合理的默认超时时间（20 分钟）
- 提供用户可配置的 Timeouts 参数
- 在文档中明确说明可能存在的延迟

### 风险 2: 并发更新冲突

**风险**: 多个 Terraform 进程同时更新同一个 DNS 记录可能导致冲突

**缓解措施**:
- 在 Update 操作中使用乐观锁机制（通过 ModifiedOn 时间戳）
- 提供明确的错误提示，建议用户重试
- 在文档中说明并发更新的限制

### 风险 3: 批量删除接口的限制

**风险**: DeleteDnsRecords API 支持批量删除，但 Terraform 的 Delete 操作是针对单个资源

**缓解措施**:
- 在 Delete 操作中只传递单个 record_id
- 利用批量删除 API 的能力，但逻辑上保持单个删除的语义

### 权衡 1: 复杂度 vs 功能完整性

**权衡**: 是否在 Provider 层面进行复杂的参数验证

**决策**: 优先保证功能完整性，依赖云 API 的参数校验

**理由**:
- 云 API 的参数校验更准确和及时
- 减少维护成本，避免与云 API 的参数规则不一致
- 用户可以通过错误信息了解问题

### 权衡 2: 性能 vs 一致性

**权衡**: 在 Read 操作中是否缓存数据以减少 API 调用

**决策**: 不缓存，每次都调用云 API 获取最新数据

**理由**:
- 保证数据的一致性
- Terraform 的 Refresh 操作频率较低，性能影响可控
- 与现有 TEO 资源保持一致

## Migration Plan

### 部署步骤

1. **代码实现**
   - 创建 resource_tc_teo_dns_record_v10.go 文件
   - 在 service_tencentcloud_teo.go 中添加服务方法
   - 创建单元测试文件 resource_tc_teo_dns_record_v10_test.go

2. **文档生成**
   - 创建 website/docs/r/teo_dns_record_v10.md 文档
   - 添加资源使用示例

3. **注册资源**
   - 在 tencentcloud/services/teo/provider.go 中注册新资源

4. **测试验证**
   - 运行单元测试确保功能正确
   - 在测试环境进行集成测试
   - 验证所有 CRUD 操作和错误场景

### 回滚策略

如果新资源存在问题，可以：
1. 从 provider.go 中移除资源注册
2. 重新构建并发布 Provider
3. 用户可以继续使用之前的版本
4. 已创建的资源不会被影响（因为资源文件仍在代码库中）

## Open Questions

无

## Implementation Notes

### 文件清单

1. `tencentcloud/services/teo/resource_tc_teo_dns_record_v10.go` - 资源实现
2. `tencentcloud/services/teo/resource_tc_teo_dns_record_v10_test.go` - 单元测试
3. `tencentcloud/services/teo/service_tencentcloud_teo.go` - 服务层方法（修改）
4. `website/docs/r/teo_dns_record_v10.md` - 资源文档
5. `tencentcloud/services/teo/provider.go` - 资源注册（修改）

### API 调用流程

1. **Create**: CreateDnsRecord → 轮询 DescribeDnsRecords 直到记录生效
2. **Read**: DescribeDnsRecords（通过 record_id 过滤）
3. **Update**: DescribeDnsRecords → ModifyDnsRecords → 轮询 DescribeDnsRecords
4. **Delete**: DeleteDnsRecords

### 关键代码模式

- 使用 `defer tccommon.LogElapsed()` 记录执行时间
- 使用 `defer tccommon.InconsistentCheck()` 进行一致性检查
- 使用 `resource.Retry()` 处理可重试错误
- 使用 `helper.String()`、`helper.IntUint64()` 等辅助函数创建 API 参数
- 使用 `helper.InterfacesHeadMap()` 处理嵌套的 Map 类型参数
