## Context

Terraform Provider for TencentCloud 目前不支持 TEO (边缘安全加速平台) 的 DNS 记录管理。TEO SDK 提供了完整的 DNS 记录 CRUD 接口，包括 CreateDnsRecord、DescribeDnsRecords、ModifyDnsRecords 和 DeleteDnsRecords。

当前 TEO 服务在 Provider 中已有其他资源实现，需要遵循现有的代码模式和结构。资源类型为 RESOURCE_KIND_GENERAL，需要实现完整的 CRUD 生命周期管理。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_10` 资源，支持完整的 CRUD 操作
- 支持 A、AAAA、MX、CNAME、TXT、NS、CAA、SRV 等 DNS 记录类型
- 支持 DNS 记录的属性：名称、类型、内容、线路、TTL、权重、优先级
- 实现正确的状态管理和幂等性
- 添加单元测试，使用 mock 的方式测试业务逻辑

**Non-Goals:**
- 不实现 DNS 记录状态修改（使用 ModifyDnsRecordsStatus 接口）
- 不实现批量操作（虽然底层 API 支持批量，但 Terraform 资源是单个资源级别）
- 不实现 DNS 记录的导入功能（可以在后续版本中添加）

## Decisions

### 1. 资源 ID 格式

**决策**: 使用复合 ID 格式 `zoneId#recordId`

**理由**:
- ZoneId 是必需的查询参数，用于唯一标识站点
- RecordId 是 DNS 记录的唯一标识
- 使用 `#` 分隔符与其他 TEO 资源保持一致
- 便于从 ID 中提取 ZoneId 用于 API 调用

**考虑的替代方案**:
- 仅使用 RecordId: 无法直接获取 ZoneId，需要额外查询
- 使用其他分隔符: 不符合 TEO 资源的现有模式

### 2. 更新操作实现

**决策**: 在 Update 函数中调用 ModifyDnsRecords API，传入单个 DNS 记录

**理由**:
- ModifyDnsRecords API 支持批量修改，传入单个记录符合 API 设计
- 避免复杂的批量逻辑，保持单个资源的语义清晰
- 减少与其他资源的耦合

**考虑的替代方案**:
- 删除后重新创建: 破坏状态，可能导致服务中断
- 使用批量修改逻辑: 增加复杂度，不需要

### 3. 幂等性实现

**决策**: 在 Create 函数中先调用 DescribeDnsRecords 检查记录是否已存在

**理由**:
- DNS 记录可能在 Provider 外部被创建
- 避免重复创建相同名称和内容的记录
- 符合 Terraform 资源的期望行为

**考虑的替代方案**:
- 仅依赖 API 的幂等性: 无法处理外部创建的资源
- 不检查: 可能创建重复记录

### 4. 可选参数处理

**决策**: Location, TTL, Weight, Priority 设置为可选参数，使用 Computed

**理由**:
- 这些参数在不同记录类型下有不同的适用性
- 云 API 提供默认值，用户可以不指定
- 使用 Computed 允许用户在配置中省略，同时在状态中显示实际值

**考虑的替代方案**:
- 设置为必填: 限制用户使用，不符合 API 设计
- 使用默认值: 与云 API 的默认值可能不一致

### 5. 轮询策略

**决策**: 在 Create 和 Update 操作后，使用 Read 函数轮询直到记录可用

**理由**:
- DNS 记录创建和修改是异步操作
- 需要确保记录生效后再返回
- 使用现有的 helper.Retry 机制，保持代码一致性

**考虑的替代方案**:
- 不等待: 可能导致后续操作失败
- 固定延迟: 不可靠，实际生效时间不确定

### 6. 测试策略

**决策**: 使用 mock 方式编写单元测试，不依赖真实云 API

**理由**:
- 避免在测试中调用真实 API，确保测试的独立性和速度
- Mock 云 API 调用，专注测试业务逻辑
- 符合项目中其他资源的测试模式

**考虑的替代方案**:
- 集成测试: 依赖真实环境和凭证，测试不稳定
- 不写测试: 不符合项目质量要求

## Risks / Trade-offs

### 风险 1: 轮询超时
**风险**: DNS 记录异步生效时间较长，可能导致轮询超时

**缓解措施**:
- 使用合理的轮询间隔和超时时间
- 在 Timeouts 配置中允许用户自定义超时时间
- 在文档中说明异步操作的特性

### 风险 2: 幂等性冲突
**风险**: 如果存在多个相同名称和内容的 DNS 记录，幂等性检查可能失败

**缓解措施**:
- 在 Create 时检查 ZoneId + Name + Type + Content 的组合唯一性
- 如果存在多个匹配记录，返回明确的错误信息
- 在文档中说明幂等性检查的限制

### 权衡 1: 批量操作 vs 单个资源
**权衡**: 底层 API 支持批量操作，但 Terraform 资源是单个资源级别

**决策**: 实现为单个资源，保持 Terraform 资源的语义清晰

### 权衡 2: 复杂性 vs 功能完整性
**权衡**: 完整实现所有 DNS 记录类型的特定参数会增加复杂度

**决策**: 仅实现通用参数，不同记录类型的特定参数通过 Content 字段传递，保持实现简洁

## Migration Plan

### 部署步骤
1. 创建资源文件 `resource_tc_teo_dns_record_10.go`
2. 在 `service_tencentcloud_teo.go` 中注册新资源
3. 创建单元测试文件 `resource_tc_teo_dns_record_10_test.go`
4. 创建文档文件 `website/docs/r/teo_dns_record_10.md`
5. 运行测试确保功能正确

### 回滚策略
- 如果发现严重问题，可以从服务注册中移除该资源
- 不影响现有资源和用户配置
- 版本回退是安全的

## Open Questions

1. 是否需要实现 DNS 记录的导入功能？
   - 当前实现不包含导入功能，可以在后续版本中根据用户需求添加

2. 是否需要支持 DNS 记录的状态修改（启用/停用）？
   - 云 API 提供了 ModifyDnsRecordsStatus 接口，但本次实现不包含此功能

3. 对于 MX 记录的 Priority 参数，是否需要验证范围？
   - 根据云 API 文档，范围是 0~50，在 Schema 中进行验证

4. 对于 Weight 参数，值为 0 时表示不解析，需要在 Update 中支持此行为吗？
   - 是的，完全按照云 API 的行为实现
