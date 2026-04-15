## Context

当前 Terraform Provider for TencentCloud 已支持 teo（EdgeOne）产品的多个资源，但尚未提供 DNS 记录管理能力。用户需要通过云 API 或控制台手动管理 DNS 记录，无法实现基础设施即代码（IaC）的标准化管理。

teo SDK 提供了完整的 DNS 记录管理接口：
- CreateDnsRecord：创建单条 DNS 记录
- DescribeDnsRecords：查询 DNS 记录列表，支持多种过滤和排序方式
- ModifyDnsRecords：批量修改 DNS 记录（最多 100 条）
- DeleteDnsRecords：批量删除 DNS 记录（最多 1000 条）

根据 vendor 目录下的 teo SDK 接口定义，DNS 记录资源需要管理以下核心字段：
- 资源标识符：RecordId（创建时生成）
- 必填字段：ZoneId（站点 ID）、Name（记录名）、Type（记录类型）、Content（记录内容）
- 可选字段：Location（解析线路）、TTL（缓存时间）、Weight（权重）、Priority（优先级）
- 只读字段：Status（状态）、CreatedOn（创建时间）、ModifiedOn（修改时间）

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_dns_record_v9` 资源的完整 CRUD 操作
- 支持所有 DNS 记录类型（A、AAAA、MX、CNAME、TXT、NS、CAA、SRV）
- 提供正确的参数校验和错误处理
- 支持异步操作的轮询机制（如果接口为异步）
- 遵循 Terraform Provider 最佳实践和项目规范
- 提供完整的单元测试和资源样例文档

**Non-Goals:**
- 批量创建/修改/删除操作（API 支持批量，但 Terraform 资粒度为单条记录）
- DNS 记录状态切换（启用/停用），该功能需要使用 ModifyDnsRecordsStatus 接口，不在本次实现范围内
- 与其他 teo 资源的关联管理
- DNS 记录的监控和告警功能

## Decisions

### 1. 资源 ID 格式

**决策**：使用复合 ID 格式 `zoneId#recordId`

**理由**：
- ZoneId 是必填字段，标识站点
- RecordId 是资源唯一标识，在创建时生成
- 复合格式便于资源读取和操作

**备选方案**：仅使用 RecordId
**拒绝理由**：读取操作需要 ZoneId 才能查询，仅 RecordId 无法获取完整信息

### 2. 批量操作 API 的使用方式

**决策**：虽然 ModifyDnsRecords 和 DeleteDnsRecords 支持批量操作，但单个 Terraform 资源实例仅管理一条 DNS 记录

**理由**：
- Terraform 资源的粒度应为单个实体，便于状态管理和依赖关系
- 用户可以通过多个资源实例实现批量管理
- 单资源单记录是 Terraform 的常见模式

### 3. 只读字段的处理

**决策**：Status、CreatedOn、ModifiedOn 作为 Computed 字段，仅在 Read 操作中填充

**理由**：
- 这些字段由云服务端生成和更新，不能由用户设置
- 根据 SDK 文档，ModifyDnsRecords 接口会忽略这些字段
- Computed 字段可以用于查看资源状态，但不参与创建/更新

### 4. 异步操作的处理

**决策**：根据 SDK 文档，CreateDnsRecord 接口返回 RecordId，但 DNS 记录的生效可能有延迟，需要在 Update 和 Delete 操作后调用 Read 接口轮询直到记录生效

**理由**：
- DNS 记录的解析生效需要时间，直接返回可能导致后续操作失败
- Terraform 的 State 需要反映资源的最新状态
- Read 操作通过 DescribeDnsRecords 接口查询，可以准确判断记录是否已生效

### 5. 记录类型相关字段的校验

**决策**：在 Schema 中为 Type 相关的字段添加 DiffSuppressFunc 和 ValidateFunc，确保字段只在对应的记录类型下生效

**理由**：
- Location 和 Weight 仅适用于 A、AAAA、CNAME 记录类型
- Priority 仅适用于 MX 记录类型
- 避免用户在不适用的记录类型下配置这些字段

### 6. 轮询重试机制

**决策**：使用 Terraform Plugin SDK v2 提供的 retry 状态刷新机制，结合 helper.Retry() 函数实现最终一致性重试

**理由**：
- 遵循项目中已有的最终一致性处理模式
- Terraform SDK 提供了完善的状态刷新机制
- 避免硬编码的轮询逻辑，提高代码可维护性

## Risks / Trade-offs

### 风险 1：DNS 记录生效延迟可能导致测试不稳定

**风险**：DNS 记录的解析生效需要时间，可能导致测试过程中出现不一致的状态

**缓解措施**：
- 在测试中使用较长的重试间隔和最大重试次数
- 在 Read 操作中实现幂等性，即使记录未完全生效也能返回当前状态
- 使用 mock 测试避免对真实云 API 的依赖

### 风险 2：批量删除接口的安全问题

**风险**：DeleteDnsRecords 接口接受 RecordIds 数组，如果实现不当可能导致误删其他记录

**缓解措施**：
- 严格校验资源 ID，确保仅删除当前资源对应的记录
- 在 Delete 操作前先调用 Read 操作确认记录存在
- 在测试中验证删除操作的正确性

### 风险 3：修改操作的字段完整性

**风险**：ModifyDnsRecords 接口需要传入完整的 DnsRecord 对象，如果字段处理不当可能导致配置丢失

**缓解措施**：
- 在 Update 操作中先 Read 当前记录状态，合并用户修改的字段
- 确保 Computed 字段（Status、CreatedOn、ModifiedOn）不参与 Update 操作
- 在测试中验证部分字段更新不会影响其他字段

### 权衡 1：代码复杂度 vs 功能完整性

**权衡**：记录类型相关的字段校验会增加代码复杂度，但能提高用户体验

**决策**：优先保证功能完整性和用户体验，接受适度的代码复杂度

**理由**：用户可能在错误的记录类型下配置不相关的字段，导致 API 调用失败

### 权衡 2：测试覆盖率 vs 开发效率

**权衡**：完整的单元测试需要 mock 所有云 API 调用，开发成本较高

**决策**：优先保证核心 CRUD 逻辑的测试覆盖，边界条件和错误处理的测试可以适当简化

**理由**：单元测试的主要目的是验证业务逻辑的正确性，而不是覆盖所有 API 场景
